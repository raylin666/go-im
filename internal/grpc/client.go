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

var _ GrpcClient = (*Client)(nil)

type GrpcClient interface {
	Account() accountPb.ServiceClient
}

type Client struct {
	ctx         context.Context
	environment system.Environment
	logger      *logger.Logger

	connects []*grpc.ClientConn

	accountClient accountPb.ServiceClient
}

func NewGrpcClient(tools *app.Tools) (grpcClient GrpcClient, cleanup func(), err error) {
	var client = &Client{
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

func (client *Client) connect() error {
	// 帐号服务客户端
	accountEndpoint := client.getAccountEndpoint()
	accountClientConn, err := dial(client.ctx, accountEndpoint, client.logger)
	if err != nil {
		client.logger.UseGrpc(client.ctx).Error(fmt.Sprintf("The account service client `%s` connected error.", accountEndpoint), zap.Error(err))
		return err
	}
	client.connects = append(client.connects, accountClientConn)
	client.accountClient = accountPb.NewServiceClient(accountClientConn)
	client.logger.UseGrpc(client.ctx).Info(fmt.Sprintf("The account service client `%s` connected successfully.", accountEndpoint))

	return nil
}

func (client *Client) close() {
	for _, conn := range client.connects {
		conn.Close()
	}
}

// getAccountEndpoint 获取帐号服务地址
func (client *Client) getAccountEndpoint() string {
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

func (client *Client) Account() accountPb.ServiceClient {
	//TODO implement me

	return client.accountClient
}
