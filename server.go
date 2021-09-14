package cache

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"github.com/wuyan94zl/go-cache/proto"
)

type Cache struct {
}

func (c *Cache) Get(ctx context.Context, req *cachepb.Request) (*cachepb.Response, error) {
	v, err := Instance.OnlyGet(req.Key)
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
