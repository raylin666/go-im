package websocket

import (
	"context"
	"go.uber.org/zap"
	"mt/config"
	"mt/pkg/repositories"
)

func ConfigServer() (config *config.Server) {
	config = ManagerInstance().cServer
	return
}

func DbRepo() (repo repositories.DbRepo) {
	repo = ManagerInstance().resource.DbRepo
	return
}

func RedisRepo() (repo repositories.RedisRepo) {
	repo = ManagerInstance().resource.RedisRepo
	return
}

func Logger(ctx context.Context) (logger *zap.Logger) {
	logger = ManagerInstance().tools.Logger().UseWebSocket(ctx)
	return
}
