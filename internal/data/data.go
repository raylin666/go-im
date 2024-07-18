package data

import (
	"context"
	"fmt"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"math/rand"
	"mt/config"
	"mt/internal/app"
	"mt/internal/grpc"
	"mt/internal/lib"
	"mt/pkg/logger"
	"mt/pkg/repositories"
	"mt/pkg/utils"
	"sync"
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
	go startTicketGrpcClientPool(ctx, srvRegister, tools.Logger())

	return &Data{
		DbRepo:    repo.DbRepo(),
		RedisRepo: repo.RedisRepo(),
	}, cleanup, nil
}

// startTicketGrpcClientPool 启动 GRPC 客户端连接池注册/销毁
func startTicketGrpcClientPool(ctx context.Context, srvRegister *lib.SrvRegister, logger *logger.Logger) {
	var (
		clientPoolsLock sync.RWMutex

		second = rand.Intn(5) + 2
		ticker = time.NewTicker(time.Duration(second) * time.Second)
	)

	for {
		select {
		case <-ticker.C:
			clientAddress := srvRegister.ClientAddress()

			// TODO 移除已关闭的服务连接池
			clientPoolsLock.RLock()
			for poolName, _ := range grpc.ClientPools() {
				if utils.InSliceByString(poolName, clientAddress) {
					continue
				}

				if grpc.DeleteClientPool(ctx, poolName) {
					logger.UseApp(ctx).Info(fmt.Sprintf("已成功移除已关闭的服务 `%s` GRPC 连接池", poolName))
				}
			}
			clientPoolsLock.RUnlock()

			// TODO 创建新启动的服务连接池
			for _, address := range clientAddress {
				if address == app.ServerIp {
					continue
				}

				if grpc.CreateClientPool(ctx, address, kratosGrpc.WithEndpoint(address)) {
					logger.UseApp(ctx).Info(fmt.Sprintf("已成功启动服务 `%s` GRPC 连接池", address))
				}
			}
		}
	}
}
