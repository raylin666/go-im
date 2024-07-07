package grpc

import (
	"context"
	"google.golang.org/grpc"
	accountPb "mt/api/v1/account"
)

type MessageClient struct {
	client accountPb.ServiceClient
}

func NewMessageClient(conn *grpc.ClientConn) *MessageClient {
	return &MessageClient{accountPb.NewServiceClient(conn)}
}

func (c *MessageClient) GenerateToken(ctx context.Context, accountId string) (string, error) {
	reply, err := c.client.GenerateToken(ctx, &accountPb.GenerateTokenRequest{AccountId: accountId})
	if err != nil {
		return "", err
	}

	return reply.Token, nil
}
