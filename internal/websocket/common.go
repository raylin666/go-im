package websocket

import (
	"context"
	"go.uber.org/zap"
	"mt/pkg/repositories"
)

const (
	codeStatusOk  uint32 = 200
	codeMessageOk string = "OK"
)

func Logger(ctx context.Context) (logger *zap.Logger) {
	logger = ManagerInstance().Logger().UseWebSocket(ctx)
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
