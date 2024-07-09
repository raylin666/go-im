package grpc

import (
	"context"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	pool "github.com/jolestar/go-commons-pool/v2"
	"sync"
)

/*
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
*/
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

	factory := pool.NewPooledObjectFactorySimple(func(ctx context.Context) (interface{}, error) {
		return kratosGrpc.Dial(ctx, opts...)
	})

	clientPoolsLock.Lock()
	defer clientPoolsLock.Unlock()

	clientPools[name] = pool.NewObjectPool(ctx, factory, &pool.ObjectPoolConfig{
		LIFO:                     false,
		MaxTotal:                 0,
		MaxIdle:                  0,
		MinIdle:                  0,
		TestOnCreate:             false,
		TestOnBorrow:             false,
		TestOnReturn:             false,
		TestWhileIdle:            false,
		BlockWhenExhausted:       false,
		MinEvictableIdleTime:     0,
		SoftMinEvictableIdleTime: 0,
		NumTestsPerEvictionRun:   0,
		EvictionPolicyName:       "",
		TimeBetweenEvictionRuns:  0,
		EvictionContext:          nil,
	})

	return clientPools[name], nil
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
