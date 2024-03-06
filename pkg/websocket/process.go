package websocket

import (
	"context"
	"go.uber.org/zap"
	"mt/pkg/logger"
	"runtime/debug"
)

type Process struct {
	ctx    context.Context
	logger *logger.Logger
}

func NewProcess(ctx context.Context, logger *logger.Logger) (process *Process) {
	process = &Process{
		ctx:    ctx,
		logger: logger,
	}

	return
}

func (p *Process) HandlerMessage(client *Client, message []byte) {
	p.logger.UseApp(p.ctx).Info("消息处理", zap.String("address", client.Addr), zap.String("message", string(message)))

	defer func() {
		if r := recover(); r != nil {
			p.logger.UseApp(p.ctx).Info("消息处理 Recover 失败", zap.Stack(string(debug.Stack())), zap.Any("recover", r))
		}
	}()
}
