package websocket

import (
	"context"
)

type EventDisposeFunc func(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

type Events struct {
	Registers map[string]EventDisposeFunc
}

func NewEvents() (events *Events) {
	events = &Events{}
	events.Registers = make(map[string]EventDisposeFunc)

	// 注册处理事件
	events.Registers["ping"] = events.Ping

	return
}

func (event *Events) Ping(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	return 200, "PONG", "{}"
}
