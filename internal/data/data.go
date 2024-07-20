package data

import (
	"context"
	"fmt"
	"math/rand"
	v1 "mt/api/v1"
	"mt/config"
	"mt/internal/app"
	"mt/internal/grpc"
	"mt/internal/lib"
	"mt/pkg/logger"
	"mt/pkg/repositories"
	"mt/pkg/utils"
	"net"
	"strings"
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
func NewData(cServer *config.Server, repo repositories.DataRepo, tools *app.Tools) (*Data, func(), error) {
	var (
		ctx = context.Background()

		// TODO 注意!!! addr 必须配置为 IP:PORT 或 :PORT 模式, 不能代理域名, 否则服务注册也会出现问题
		addr = cServer.Grpc.GetAddr()
	)

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("请务必将当前 GRPC 地址 `%s` 配置为 IP:PORT 或 :PORT 模式, 否则服务注册将出现问题, 导致服务错误", addr))
	}

	// .env.yaml 配置的 addr IP一般是 127.0.0.1, 需要替换为当前服务的IP地址后注册服务
	addr = fmt.Sprintf("%s:%d", app.LocalServerIp, tcpAddr.Port)

	srvRegister := lib.NewSrvRegister(ctx, addr, repo, tools)

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
	go startTicketGrpcClientPool(ctx, addr, srvRegister, tools.Logger())

	return &Data{
		DbRepo:    repo.DbRepo(),
		RedisRepo: repo.RedisRepo(),
	}, cleanup, nil
}

// startTicketGrpcClientPool 启动 GRPC 客户端连接池注册/销毁
func startTicketGrpcClientPool(ctx context.Context, addr string, srvRegister *lib.SrvRegister, logger *logger.Logger) {
	var (
		clientPoolsLock sync.RWMutex
		connErrorLock   sync.RWMutex

		connError = make(map[string]int)
		second    = rand.Intn(5) + 2
		ticker    = time.NewTicker(time.Duration(second) * time.Second)

		localAddr = "127.0.0.1"
		parseAddr = strings.Split(addr, ":")
	)

	for {
		select {
		case <-ticker.C:
			clientAddress := srvRegister.ClientAddress()

			// TODO 移除已关闭的服务连接池
			clientPoolsLock.RLock()
			for address, _ := range grpc.ClientPools() {
				parseAddress := strings.Split(address, ":")
				// 本机地址无需处理
				if parseAddress[0] == localAddr {
					continue
				}

				if utils.InSliceByString(address, clientAddress) {
					// 检测该服务器的连接是否可用, 当连续检测超过10次失败时, 将服务移除出注册
					conn, err := grpc.Dial(ctx, address)
					if err == nil {
						heartbeatClient := v1.NewHeartbeatClient(conn)
						_, err = heartbeatClient.PONE(ctx, nil)
					}

					if err != nil {
						if _, ok := err.(net.Error); ok {
							// 网络错误, 无法调通
							connErrorLock.Lock()
							connError[address] += 1
							connErrorLock.Unlock()

							logger.UseApp(ctx).Warn(fmt.Sprintf("检测 GRPC 服务 `%s` 连通性失败, 无法连接服务器, 连续失败次数为 `%d` 次", address, connError[address]))

							connErrorLock.RLock()
							if connError[address] >= 10 {
								srvRegister.UnRegister()
							}
							connErrorLock.RUnlock()
						} else {
							// 调用错误, 比如方法不存在等 ...

							connError[address] = 0
						}
					} else {
						// GRPC 服务器连接成功, 重置连续失败次数
						connError[address] = 0
					}

					conn.Close()

					continue
				}

				if grpc.DeleteClientPool(ctx, address) {
					logger.UseApp(ctx).Info(fmt.Sprintf("已成功移除已关闭的服务 `%s` GRPC 连接池", address))
				}
			}
			clientPoolsLock.RUnlock()

			// TODO 创建新启动的服务连接池
			for _, address := range clientAddress {
				// 当前服务的GRPC连接池也需要创建, 因为需要用到请求API
				// 当创建的连接池为当前服务时, 将IP地址转换为 127.0.0.1
				// 理论上 address 都是内网地址, 为避免其他不可测性, 故将当前服务的连接地址直接转换为本机更为稳妥
				parseAddress := strings.Split(address, ":")
				if parseAddress[0] == app.LocalServerIp && parseAddress[1] == parseAddr[1] {
					address = fmt.Sprintf("%s:%s", localAddr, parseAddress[1])
				}

				if grpc.CreateClientPool(ctx, address) {
					logger.UseApp(ctx).Info(fmt.Sprintf("已成功启动服务 `%s` GRPC 连接池", address))
				}
			}
		}
	}
}
