package cache

import (
	"context"
	"fmt"
	"github.com/wuyan94zl/go-cache/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Cache struct {
}

func (c *Cache) Set(ctx context.Context, req *cachepb.SetRequest) (*cachepb.Response, error) {
	v, err := Instance.onlySet(req.Key, []byte(req.Value), req.Ttl)
	if err != nil {
		return &cachepb.Response{}, err
	} else {
		return &cachepb.Response{Value: v.B}, err
	}
}

func (c *Cache) SetNX(ctx context.Context, req *cachepb.SetRequest) (*cachepb.SetResponse, error) {
	if Instance.onlySexNX(req.Key, []byte(req.Value), req.Ttl) {
		return &cachepb.SetResponse{Value: true}, nil
	} else {
		return &cachepb.SetResponse{Value: false}, fmt.Errorf("key is exist")
	}
}

func (c *Cache) Get(ctx context.Context, req *cachepb.Request) (*cachepb.Response, error) {
	v, err := Instance.onlyGet(req.Key)
	if err != nil {
		return &cachepb.Response{}, err
	} else {
		return &cachepb.Response{Value: v.B}, err
	}
}

func listen(port string) {
	addr := "0.0.0.0:" + port
	listen, _ := net.Listen("tcp", addr)
	// 实例化grpc Server
	s := grpc.NewServer()
	// 注册OrderServer
	cachepb.RegisterGroupCacheServer(s, &Cache{})
	log.Println("Listen on " + addr)
	go s.Serve(listen)
	/** end grpc 服务启动 */
}
