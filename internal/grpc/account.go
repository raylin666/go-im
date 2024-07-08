package grpc

import (
	"context"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	accountPb "mt/api/v1/account"
)

type AccountClient struct {
	client accountPb.ServiceClient
}

func NewAccountClient(ctx context.Context, endpoint string, conn *grpc.ClientConn) *AccountClient {
	clientPool, err := GetClientPool(ctx, endpoint, kratosGrpc.WithEndpoint(endpoint))
	if err != nil {
		return nil
	}

	return &AccountClient{accountPb.NewServiceClient(clientPool.Get())}
}

func (c *AccountClient) GenerateToken(ctx context.Context, accountId string) (string, error) {
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
