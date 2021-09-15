package cache

import (
	"context"
	"fmt"
	"github.com/wuyan94zl/go-cache/byteview"
	"github.com/wuyan94zl/go-cache/consistenthash"
	"github.com/wuyan94zl/go-cache/lru"
	"github.com/wuyan94zl/go-cache/proto"
	"github.com/wuyan94zl/go-cache/singleflight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"sync"
)

type cache struct {
	mu  sync.RWMutex
	lru *lru.Cache
}

func (c *cache) add(key string, value byteview.ByteView, ttl int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		panic("lru is nil")
	}
	c.lru.Add(key, value, ttl)
}

func (c *cache) get(key string) (value byteview.ByteView, ok bool) {
	c.mu.Lock()
	c.mu.Unlock()
	if c.lru == nil {
		panic("lru is nil")
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(byteview.ByteView), ok
	}
	return
}

func (c *cache) setNx(key string, value byteview.ByteView, ttl int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		panic("lru is nil")
	}
	if _, ok := c.lru.Get(key); ok {
		return false
	} else {
		c.lru.Add(key, value, ttl)
		return true
	}
}

type CallBack interface {
	Get(key string, params map[string]interface{}) ([]byte, error)
}

type CallBackFunc func(key string, params map[string]interface{}) ([]byte, error)

func (f CallBackFunc) Get(key string, params map[string]interface{}) ([]byte, error) {
	return f(key, params)
}

type Group struct {
	callback  CallBack
	mainCache cache
	isAdd     bool
	loader    *singleflight.Group
	hash      *consistenthash.Map
	self      string
}

var Instance *Group
var maxLen = 10000
var backupInterval = 60

type GrpcConfig struct {
	Port        string
	CurrentHost string
	AllHosts    []string
}

type Config struct {
	MaxLen         int
	BackupInterval int
	Grpc           *GrpcConfig
}

func Default(config *Config) {
	if Instance == nil {
		if config != nil && config.MaxLen > 0 {
			maxLen = config.MaxLen
		}
		if config != nil && config.BackupInterval > 0 {
			backupInterval = config.BackupInterval
		}
		Instance = &Group{
			mainCache: cache{lru: lru.New(maxLen, backupInterval, nil)},
			loader:    singleflight.NewGroup(),
			hash:      consistenthash.New(3, nil),
		}
		Instance.init(config)
	}
}

func (g *Group) init(c *Config) {
	if c != nil && c.Grpc != nil {
		g.self = c.Grpc.CurrentHost
		g.hash.Add(c.Grpc.AllHosts...)
		listen(c.Grpc.Port)
	}
}

func (g *Group) Set(key string, value string, ttl int64) bool {
	if url, ok := g.isLocal(key); !ok {
		return g.grpcSetData(url, key, []byte(value), ttl)
	} else {
		g.mainCache.add(key, byteview.ByteView{B: []byte(value)}, ttl)
		return true
	}
}

func (g *Group) Get(key string) ([]byte, error) {
	var bytes []byte
	var err error
	if url, ok := g.isLocal(key); !ok {
		if bytes, err = g.grpcGetData(url, key); err == nil {
			return bytes, err
		} else {
			return nil, err
		}
	} else {
		val, err := g.mainCache.get(key)
		if err {
			return val.B, nil
		} else {
			return nil, fmt.Errorf("key is not exist")
		}
	}
}

func (g *Group) SetNX(key string, value string, ttl int64) bool {
	if url, ok := g.isLocal(key); !ok {
		return g.grpcSetNXData(url, key, []byte(value), ttl)
	} else {
		return g.mainCache.setNx(key, byteview.ByteView{B: []byte(value)}, ttl)
	}
}

func (g *Group) CallBackFunc(callback CallBack) *Group {
	g.callback = callback
	return g
}

func (g *Group) Cache(key string, params map[string]interface{}, ttl int64) (byteview.ByteView, error) {
	if key == "" {
		return byteview.ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		return v, nil
	}
	return g.load(key, params, ttl)
}

func (g *Group) load(key string, params map[string]interface{}, ttl int64) (byteview.ByteView, error) {
	v, err := g.loader.Do(key, func() (interface{}, error) {
		return g.getLocally(key, params, ttl)
	})
	return v.(byteview.ByteView), err
}

func (g *Group) getLocally(key string, params map[string]interface{}, ttl int64) (byteview.ByteView, error) {
	var bytes []byte
	var err error

	if url, ok := g.isLocal(key); !ok {
		if bytes, err = g.grpcGetData(url, key); err == nil {
			return byteview.ByteView{B: bytes}, nil
		}
	}

	bytes, err = g.callback.Get(key, params)
	if err != nil {
		return byteview.ByteView{}, err
	}
	value := byteview.ByteView{B: bytes}
	if !g.isAdd {
		ttl = 10
	}
	g.populateCache(key, value, ttl)
	return value, nil
}

func (g *Group) isLocal(key string) (string, bool) {
	g.isAdd = true
	url := g.hash.Get(key)
	if url == g.self {
		return "", true
	}
	g.isAdd = false
	return url, false
}

func (g *Group) grpcGetData(url string, key string) ([]byte, error) {
	conn, err := grpc.Dial(url, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("连接服务器失败:%v", err)
	}
	cli := NewClient(conn)
	res, err := cli.Get(context.Background(), &cachepb.Request{Key: key})
	if err != nil {
		return nil, err
	}
	return res.Value, nil
}

func (g *Group) grpcSetData(url string, key string, value []byte, ttl int64) bool {
	conn, err := grpc.Dial(url, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithInsecure())
	if err != nil {
		return false
	}
	cli := NewClient(conn)
	_, err = cli.Set(context.Background(), &cachepb.SetRequest{Key: key, Value: string(value), Ttl: ttl})
	if err != nil {
		return false
	}
	return true
}

func (g *Group) grpcSetNXData(url string, key string, value []byte, ttl int64) bool {
	conn, err := grpc.Dial(url, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithInsecure())
	if err != nil {
		return false
	}
	cli := NewClient(conn)
	_, err = cli.SetNX(context.Background(), &cachepb.SetRequest{Key: key, Value: string(value), Ttl: ttl})
	if err != nil {
		return false
	}
	return true
}

func (g *Group) populateCache(key string, value byteview.ByteView, ttl int64) {
	g.mainCache.add(key, value, ttl)
}

func (g *Group) OnlySet(key string, value []byte, ttl int64) (byteview.ByteView, error) {
	v := byteview.ByteView{B: value}
	g.mainCache.add(key, v, ttl)
	return v, nil
}

func (g *Group) OnlyGet(key string) (byteview.ByteView, error) {
	if key == "" {
		return byteview.ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		return v, nil
	} else {
		return byteview.ByteView{}, fmt.Errorf("key is not exist")
	}
}

func (g *Group) OnlySexNX(key string, value []byte, ttl int64) bool {
	v := byteview.ByteView{B: value}
	return g.mainCache.setNx(key, v, ttl)
}
