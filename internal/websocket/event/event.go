package event

import (
	"github.com/google/wire"
	"mt/internal/websocket"
)

// ProviderSet is events providers.
var ProviderSet = wire.NewSet(NewEvents)

type events struct {
	relation map[string]websocket.EventDisposeFunc
}

// NewEvents 消息事件
func NewEvents() websocket.Events {
	return &events{}
}

// GetAll 获取所有事件对应处理器
func (event *events) GetAll() map[string]websocket.EventDisposeFunc {
	event.relation = make(map[string]websocket.EventDisposeFunc)

	// Ping 心跳检测处理器
	event.relation[websocket.EventPing] = event.Ping

	return event.relation
}

// defaultEventResponse 默认事件返回值
func defaultEventResponse() (code uint32, msg string, data interface{}) {
	code, msg, data = 200, "OK", nil
	return
}
