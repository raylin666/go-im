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

func (c *AccountClient) GenerateToken(ctx context.Context, accountId string) (string, error) {
	defer c.pool.Put(c.pool.Get())
	reply, err := c.client.GenerateToken(ctx, &accountPb.GenerateTokenRequest{AccountId: accountId})
	if err != nil {
		return "", err
	}

	return reply.Token, nil
}

/*grpcConn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint("127.0.0.1:10011"))
if err != nil {
fmt.Println(err)
}

grpcClient := clientGrpc.NewAccountClient(ctx, "127.0.0.1:10011", grpcConn)
reply, err := grpcClient.GenerateToken(ctx, "91283746167")
fmt.Println(reply)
*/
