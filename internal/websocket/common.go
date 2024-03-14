package websocket

import (
	"context"
	"go.uber.org/zap"
)

const (
	codeStatusOk  = 200
	codeMessageOk = "OK"
)

func Logger(ctx context.Context) (logger *zap.Logger) {
	logger = ManagerInstance().Logger().UseWebSocket(ctx)
	return
}
