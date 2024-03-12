package websocket

import (
	"context"
	"go.uber.org/zap"
)

func Logger(ctx context.Context) (logger *zap.Logger) {
	logger = ManagerInstance().Logger().UseWebSocket(ctx)
	return
}
