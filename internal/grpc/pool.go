package grpc

import (
	"context"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	goPool "github.com/jolestar/go-commons-pool/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
)

var (
	clientPoolsLock sync.RWMutex
	clientPools     = make(map[string]*goPool.ObjectPool)
)

func ClientPools() map[string]*goPool.ObjectPool { return clientPools }

func CreateClientPool(ctx context.Context, name string, opts ...kratosGrpc.ClientOption) bool {
	if _, ok := clientPools[name]; ok {
		return false
	}

	factory := goPool.NewPooledObjectFactorySimple(func(ctx context.Context) (interface{}, error) {
		return Dial(ctx, name, opts...)
	})

	clientPoolsLock.Lock()
	defer clientPoolsLock.Unlock()

	clientPools[name] = goPool.NewObjectPool(ctx, factory, &goPool.ObjectPoolConfig{
		LIFO:                     false,
		MaxTotal:                 100,
		MaxIdle:                  10,
		MinIdle:                  5,
		TestOnCreate:             false,
		TestOnBorrow:             false,
		TestOnReturn:             false,
		TestWhileIdle:            false,
		BlockWhenExhausted:       false,
		MinEvictableIdleTime:     0,
		SoftMinEvictableIdleTime: 0,
		NumTestsPerEvictionRun:   3,
		EvictionPolicyName:       "",
		TimeBetweenEvictionRuns:  0,
		EvictionContext:          ctx,
	})

	return true
}

func DeleteClientPool(ctx context.Context, name string) bool {
	if _, ok := clientPools[name]; ok {
		clientPoolsLock.Lock()
		defer clientPoolsLock.Unlock()
		clientPools[name].Close(ctx)
		delete(clientPools, name)
		return true
	}

	return false
}

type pool struct {
	ctx  context.Context
	name string
}

type Pool interface {
	Get() *grpc.ClientConn
	Put(conn *grpc.ClientConn)
}

func NewPool(ctx context.Context, name string) Pool {
	return &pool{ctx: ctx, name: name}
}

func (pool *pool) Get() *grpc.ClientConn {
	clientPoolsLock.RLock()
	defer clientPoolsLock.RUnlock()

	object, ok := clientPools[pool.name]
	if !ok {
		return nil
	}

	connFunc := func() *grpc.ClientConn {
		objectValue, err := object.BorrowObject(pool.ctx)
		if err != nil {
			return nil
		}

		return objectValue.(*grpc.ClientConn)
	}

	conn := connFunc()

	// 如果连接关闭或失败
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		_ = conn.Close()
		_ = object.AddObject(pool.ctx)
		conn = connFunc()
	}

	return conn
}

func (pool *pool) Put(conn *grpc.ClientConn) {
	clientPoolsLock.RLock()
	defer clientPoolsLock.RUnlock()

	object, ok := clientPools[pool.name]
	if !ok {
		return
	}

	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		_ = conn.Close()
		_ = object.AddObject(pool.ctx)
		return
	}

	_ = object.ReturnObject(pool.ctx, conn)
}
