package event

import (
	"context"
	"mt/internal/websocket"
	"strings"
	"sync"
)

const (
	MessageEventPing       = "ping"
	MessageEventBind       = "bind"
	MessageEventC2CMessage = "c2c_message"
)

type MessageEventDisposeFunc func(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{})

type MessageEvent interface {
	// HasClientSupport 判断是否客户端请求所支持的消息事件, 不在指定的消息事件客户端无法调用
	HasClientSupport(event string) bool
	// GetDisposeFunc 获取消息事件的处理函数
	GetDisposeFunc(event string) (MessageEventDisposeFunc, bool)

	/**
	消息事件
	*/
	// Ping 心跳检测
	Ping(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{})
	// Bind 客户端和账号信息绑定
	Bind(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{})
	// C2CMessage 发送C2C消息
	C2CMessage(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{})
}

type messageEvent struct {
	disposeFuncs    map[string]MessageEventDisposeFunc
	disposeFuncLock sync.RWMutex
}

func NewMessageEvent() MessageEvent {
	var event = new(messageEvent)

	event.registerDisposeFunc()

	return event
}

// registerDisposeFunc 注册消息事件处理函数
func (event *messageEvent) registerDisposeFunc() {
	event.disposeFuncs = make(map[string]MessageEventDisposeFunc)

	// 心跳检测处理器
	event.disposeFuncs[MessageEventPing] = event.Ping
	// 客户端和账号信息绑定
	event.disposeFuncs[MessageEventBind] = event.Bind
	// 发送C2C消息
	event.disposeFuncs[MessageEventC2CMessage] = event.C2CMessage
}

// HasClientSupport 判断是否客户端请求所支持的消息事件, 不在指定的消息事件客户端无法调用
func (event *messageEvent) HasClientSupport(name string) bool {
	//TODO implement me

	var events = []string{
		MessageEventPing,
		MessageEventBind,
		MessageEventC2CMessage,
	}

	for _, eventValue := range events {
		if strings.Contains(eventValue, name) {
			return true
		}
	}

	return false
}

// DisposeFuncs 获取所有消息事件的处理函数
func (event *messageEvent) GetDisposeFunc(name string) (MessageEventDisposeFunc, bool) {
	//TODO implement me

	event.disposeFuncLock.RLock()
	defer event.disposeFuncLock.RUnlock()

	if f, ok := event.disposeFuncs[name]; ok {
		return f, true
	}

	return nil, false
}
