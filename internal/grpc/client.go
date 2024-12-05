package grpc

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	accountPb "mt/api/v1/account"
	"mt/internal/app"
	"mt/pkg/logger"
)

const (
	accountGrpcClientEndpoint = "127.0.0.1:10011"
)

type GrpcClient struct {
	ctx    context.Context
	logger *logger.Logger

	Connects []*grpc.ClientConn

	Account accountPb.ServiceClient
}

func NewGrpcClient(tools *app.Tools) (client *GrpcClient, cleanup func(), err error) {
	client = &GrpcClient{ctx: context.TODO(), logger: tools.Logger()}

	cleanup = func() {
		client.close()
		tools.Logger().UseGrpc(client.ctx).Info("closing the grpc clients successfully.")
	}

	err = client.connect()

	return client, cleanup, err
}

func (client *GrpcClient) connect() error {
	// 帐号服务客户端
	accountClientConn, err := dial(client.ctx, accountGrpcClientEndpoint)
	if err != nil {
		client.logger.UseGrpc(client.ctx).Error(fmt.Sprintf("The account service client `%s` connected error.", accountGrpcClientEndpoint), zap.Error(err))
		return err
	}
	client.Account = accountPb.NewServiceClient(accountClientConn)
	client.Connects = append(client.Connects, accountClientConn)
	client.logger.UseGrpc(client.ctx).Info(fmt.Sprintf("The account service client `%s` connected successfully.", accountGrpcClientEndpoint))

	return nil
}

func (client *GrpcClient) close() {
	for _, conn := range client.Connects {
		conn.Close()
	}
}
