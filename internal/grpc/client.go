package grpc

import (
	"context"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	goGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func grpcDial(ctx context.Context, address string, opts ...kratosGrpc.ClientOption) (*goGrpc.ClientConn, error) {
	// 重新序列优先级设置
	var newOpts = []kratosGrpc.ClientOption{
		kratosGrpc.WithEndpoint(address),
		kratosGrpc.WithOptions(goGrpc.WithTransportCredentials(insecure.NewCredentials())),
	}

	for _, opt := range opts {
		newOpts = append(newOpts, opt)
	}

	return kratosGrpc.Dial(ctx, newOpts...)
}
