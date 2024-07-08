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

var (
	clientPools     map[string]ClientPool
	clientPoolsLock sync.RWMutex
)

func GetClientPool(ctx context.Context, name string, opts ...kratosGrpc.ClientOption) (ClientPool, error) {
	if _, ok := clientPools[name]; ok {
		return clientPools[name], nil
	}

	pool, err := NewClientPool(ctx, opts...)
	if err != nil {
		return nil, err
	}

	clientPoolsLock.Lock()
	defer clientPoolsLock.Unlock()
	clientPools[name] = pool
	return pool, nil
}

func DelClientPool(ctx context.Context, name string) bool {
	if _, ok := clientPools[name]; ok {
		clientPoolsLock.Lock()
		defer clientPoolsLock.Unlock()
		delete(clientPools, name)
		return true
	}

	return false
}
