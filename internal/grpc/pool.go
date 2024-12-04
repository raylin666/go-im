package grpc

import (
	"context"
	"fmt"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	"sync"
)

type GRPCPool struct {
	ctx         context.Context
	clients     []*grpc.ClientConn
	clientOpts  []*kratosGrpc.ClientOption
	mutex       sync.Mutex
	address     string
	currentSize int
	maxSize     int
}

func NewGRPCPool(ctx context.Context, address string, initSize, maxSize int, opts ...kratosGrpc.ClientOption) (*GRPCPool, error) {
	pool := &GRPCPool{
		ctx:     ctx,
		address: address,
		maxSize: maxSize,
	}

	for i := 0; i < initSize; i++ {
		conn, err := grpcDial(ctx, address, opts...)
		if err != nil {
			return nil, err
		}

		pool.clients = append(pool.clients, conn)
		pool.currentSize++
	}

	return pool, nil
}

func (p *GRPCPool) Get() (*grpc.ClientConn, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.clients) == 0 && p.currentSize < p.maxSize {
		// Create new connection if pool is empty and below max size
		conn, err := grpcDial(ctx, targetAddress, opts...)
		if err != nil {
			return nil, err
		}
		p.currentSize++
		return conn, nil
	} else if len(p.clients) == 0 {
		return nil, fmt.Errorf("no available connections")
	}

	// Return a connection from the pool
	conn := p.clients[len(p.clients)-1]
	p.clients = p.clients[:len(p.clients)-1]
	return conn, nil
}

// Put returns a gRPC client to the pool.
func (p *GRPCPool) Put(conn *grpc.ClientConn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.clients = append(p.clients, conn)
}

// Close closes all the connections in the pool.
func (p *GRPCPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, conn := range p.clients {
		conn.Close()
	}
	p.clients = nil
	p.currentSize = 0
}
