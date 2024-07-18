package data

import (
	"context"
	"fmt"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"mt/config"
	"mt/internal/app"
	"mt/internal/grpc"
	"mt/internal/lib"
	"mt/pkg/repositories"
	"time"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDataRepo,
	NewHeartbeatRepo,
	NewAccountRepo)

func NewDataRepo(tools *app.Tools, data *config.Data) repositories.DataRepo {
	return repositories.NewDataRepo(tools.Logger(), data)
}

// Data .
type Data struct {
	// TODO wrapped database client
	DbRepo    repositories.DbRepo
	RedisRepo repositories.RedisRepo
}

// NewData .
func NewData(repo repositories.DataRepo, tools *app.Tools) (*Data, func(), error) {
	var ctx = context.Background()
	var srvRegister = lib.NewSrvRegister(ctx, repo, tools)

	cleanup := func() {
		// 服务下线
		srvRegister.UnRegister()

		// 资源关闭
		for dbName, dbRepo := range repo.DbRepo().All() {
			_ = dbRepo.Close()
			tools.Logger().UseApp(ctx).Info(fmt.Sprintf("closing the data resource: %s db.repo.", dbName))
		}

		for redisName, redisRepo := range repo.RedisRepo().All() {
			_ = redisRepo.Close()
			tools.Logger().UseApp(ctx).Info(fmt.Sprintf("closing the data resource: %s db.repo.", redisName))
		}
	}

	// 服务注册
	srvRegister.Register()

	// 启动 GRPC 客户端连接池注册/销毁
	go startTicketGrpcClientPool(ctx, srvRegister)

	return &Data{
		DbRepo:    repo.DbRepo(),
		RedisRepo: repo.RedisRepo(),
	}, cleanup, nil
}

// startTicketGrpcClientPool 启动 GRPC 客户端连接池注册/销毁
func startTicketGrpcClientPool(ctx context.Context, srvRegister *lib.SrvRegister) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			for _, address := range srvRegister.ClientAddress() {
				grpc.CreateClientPool(ctx, address, kratosGrpc.WithEndpoint(address))
			}
		}
	}
}
