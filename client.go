package cache

import (
	"context"
	"google.golang.org/grpc"
	"github.com/wuyan94zl/go-cache/proto"
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
