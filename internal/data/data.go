package data

import (
	"context"
	"fmt"
	"mt/config"
	"mt/internal/app"
	"mt/internal/lib"
	"mt/pkg/repositories"

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
		repo.DB(repositories.DbConnectionDefaultName).Close()
		tools.Logger().UseApp(ctx).Info(fmt.Sprintf("closing the data resource: %s db.repo.", repositories.DbConnectionDefaultName))
		repo.Redis(repositories.RedisConnectionDefaultName).Close()
		tools.Logger().UseApp(ctx).Info(fmt.Sprintf("closing the data resource: %s redis.repo.", repositories.RedisConnectionDefaultName))
	}

	// 服务注册
	srvRegister.Register()

	return &Data{
		DbRepo:    repo.DbRepo(),
		RedisRepo: repo.RedisRepo(),
	}, cleanup, nil
}
