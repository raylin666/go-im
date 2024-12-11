package grpc

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	goGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func dial(ctx context.Context, endpoint string, opts ...kratosGrpc.ClientOption) (*goGrpc.ClientConn, error) {
	// 重新序列优先级设置
	var newOpts = []kratosGrpc.ClientOption{
		kratosGrpc.WithEndpoint(endpoint),
		kratosGrpc.WithOptions(goGrpc.WithTransportCredentials(insecure.NewCredentials())),
		kratosGrpc.WithHealthCheck(true),
		kratosGrpc.WithTimeout(6 * time.Second),
		kratosGrpc.WithMiddleware(metadata.Client()),
	}

	for _, opt := range opts {
		newOpts = append(newOpts, opt)
	}

	return kratosGrpc.Dial(ctx, newOpts...)
}
