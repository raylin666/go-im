package websocket

import (
	"context"
	"go.uber.org/zap"
)

const (
	CodeStatusOk  = 200
	CodeMessageOk = "OK"
)

func Logger(ctx context.Context) (logger *zap.Logger) {
	logger = ManagerInstance().Logger().UseWebSocket(ctx)
	return
}
