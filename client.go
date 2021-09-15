package cache

import (
	"context"
	"github.com/wuyan94zl/go-cache/proto"
	"google.golang.org/grpc"
)

type defaultClient struct {
	conn *grpc.ClientConn
}

func NewClient(cli *grpc.ClientConn) *defaultClient {
	return &defaultClient{
		conn: cli,
	}
}

func (d *defaultClient) Get(ctx context.Context, request *cachepb.Request) (*cachepb.Response, error) {
	cli := cachepb.NewGroupCacheClient(d.conn)
	return cli.Get(ctx, request)
}

func (d *defaultClient) Set(ctx context.Context, request *cachepb.SetRequest) (*cachepb.Response, error) {
	cli := cachepb.NewGroupCacheClient(d.conn)
	return cli.Set(ctx, request)
}

func (d *defaultClient) SetNX(ctx context.Context, request *cachepb.SetRequest) (*cachepb.SetResponse, error) {
	cli := cachepb.NewGroupCacheClient(d.conn)
	return cli.SetNX(ctx, request)
}
