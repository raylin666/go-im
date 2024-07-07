package grpc

import (
	"context"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
)

type clientPool struct {
	pool sync.Pool
}

type ClientPool interface {
	Get() *grpc.ClientConn
	Put(conn *grpc.ClientConn)
}

func (client *clientPool) Get() *grpc.ClientConn {
	conn := client.pool.Get().(*grpc.ClientConn)
	// 如果连接关闭或失败
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		conn = client.pool.New().(*grpc.ClientConn)
	}

	return conn
}

func (client *clientPool) Put(conn *grpc.ClientConn) {
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		conn = client.pool.New().(*grpc.ClientConn)
	}

	client.pool.Put(conn)
}

func NewClientPool(ctx context.Context, opts ...kratosGrpc.ClientOption) (ClientPool, error) {
	return &clientPool{
		pool: sync.Pool{
			New: func() any {
				conn, err := kratosGrpc.Dial(ctx, opts...)
				if err != nil {
					return nil
				}

				return conn
			},
		},
	}, nil
}
