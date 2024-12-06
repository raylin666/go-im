package grpc

import (
	"context"
	"fmt"
	"github.com/raylin666/go-utils/server/system"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	accountPb "mt/api/v1/account"
	"mt/internal/app"
	"mt/pkg/logger"
)

type GrpcClient struct {
	ctx         context.Context
	environment system.Environment
	logger      *logger.Logger

	connects []*grpc.ClientConn

	Account accountPb.ServiceClient
}

func NewGrpcClient(tools *app.Tools) (client *GrpcClient, cleanup func(), err error) {
	client = &GrpcClient{
		ctx:         context.TODO(),
		environment: tools.Environment(),
		logger:      tools.Logger(),
	}

	cleanup = func() {
		client.close()
		tools.Logger().UseGrpc(client.ctx).Info("closing the grpc clients successfully.")
	}

	err = client.connect()

	return client, cleanup, err
}

func (client *GrpcClient) connect() error {
	// 帐号服务客户端
	accountEndpoint := client.getAccountEndpoint()
	accountClientConn, err := dial(client.ctx, accountEndpoint)
	if err != nil {
		client.logger.UseGrpc(client.ctx).Error(fmt.Sprintf("The account service client `%s` connected error.", accountEndpoint), zap.Error(err))
		return err
	}
	client.connects = append(client.connects, accountClientConn)
	client.Account = accountPb.NewServiceClient(accountClientConn)
	client.logger.UseGrpc(client.ctx).Info(fmt.Sprintf("The account service client `%s` connected successfully.", accountEndpoint))

	return nil
}

func (client *GrpcClient) close() {
	for _, conn := range client.connects {
		conn.Close()
	}
}

// getAccountEndpoint 获取帐号服务地址
func (client *GrpcClient) getAccountEndpoint() string {
	if client.environment.IsProd() {
		return ProdAccountGrpcClientEndpoint
	}

	if client.environment.IsPre() {
		return PreAccountGrpcClientEndpoint
	}

	if client.environment.IsTest() {
		return TestAccountGrpcClientEndpoint
	}

	return DevAccountGrpcClientEndpoint
}
