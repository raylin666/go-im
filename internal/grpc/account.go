package grpc

import (
	"context"
	accountPb "mt/api/v1/account"
)

type AccountClient struct {
	pool   Pool
	client accountPb.ServiceClient
}

func NewAccountClient(ctx context.Context, endpoint string) *AccountClient {
	var pool = NewPool(ctx, endpoint)
	return &AccountClient{pool: pool, client: accountPb.NewServiceClient(pool.Get())}
}

func (c *AccountClient) GenerateToken(ctx context.Context, accountId string) (*accountPb.GenerateTokenResponse, error) {
	defer c.pool.Put(c.pool.Get())
	return c.client.GenerateToken(ctx, &accountPb.GenerateTokenRequest{AccountId: accountId})
}
