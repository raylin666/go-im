package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/constant/defined"
	model2 "mt/internal/websocket/model"
	"mt/pkg/logger"
	"runtime/debug"
	"sync"
)

type DisposeFunc func(client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

var (
	handlers        = make(map[string]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// RegisterHandler 注册处理器
func RegisterHandler(key string, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

// getHandlers 获取所有处理器
func getHandlers(key string) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]

	return
}

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
			p.logger.UseApp(p.ctx).Error("消息处理 Recover 失败", zap.Stack(string(debug.Stack())), zap.Any("recover", r))
		}
	}()

	request := &model2.Request{}

	err := json.Unmarshal(message, request)
	if err != nil {
		p.logger.UseApp(p.ctx).Error("消息处理 JSON 解码失败", zap.ByteString("message", message), zap.Error(err))
		client.SendMessage([]byte(defined.ErrorDataValidateError.GetReason()))
		return
	}

	requestData, err := json.Marshal(request.Data)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
		client.SendMessage([]byte(defined.ErrorDataHandlerError.GetReason()))
		return
	}

	var (
		code uint32
		msg  string
		data interface{}

		seq   = request.Seq
		event = request.Event
	)

	// 采用 MAP 注册的方式
	if value, ok := getHandlers(event); ok {
		code, msg, data = value(client, seq, requestData)
	} else {
		code = 404
		fmt.Println("处理数据 路由不存在", client.Addr, "cmd", event)
	}

	msg = "我这块代码需要重新修改"

	responseHead := model2.NewResponseHead(seq, event, code, msg, data)

	headByte, err := json.Marshal(responseHead)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)

		return
	}

	client.SendMessage(headByte)

	fmt.Println("acc_response send", client.Addr, client.AppId, client.UserId, "cmd", event, "code", code)

	return
}
