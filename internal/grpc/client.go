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
	logger *logger.Logger

	Connects []*grpc.ClientConn

	Account accountPb.ServiceClient
}

func NewGrpcClient(tools *app.Tools) (grpc *GrpcClient, err error) {
	grpc = &GrpcClient{logger: tools.Logger()}
	err = grpc.connect()
	if err != nil {
		return nil, err
	}

	return grpc, nil
}

func (client *GrpcClient) connect() error {
	var ctx = context.TODO()

	// 帐号服务客户端
	accountClientConn, err := dial(ctx, accountGrpcClientEndpoint)
	if err != nil {
		client.logger.UseApp(ctx).Error(fmt.Sprintf("The account service client `%s` connected error.", accountGrpcClientEndpoint), zap.Error(err))
		return err
	}
	client.Account = accountPb.NewServiceClient(accountClientConn)
	client.Connects = append(client.Connects, accountClientConn)
	client.logger.UseApp(ctx).Info(fmt.Sprintf("The account service client `%s` connected successfully.", accountGrpcClientEndpoint))

	return nil
}

func (client *GrpcClient) close() {
	for _, conn := range client.Connects {
		conn.Close()
	}
}
