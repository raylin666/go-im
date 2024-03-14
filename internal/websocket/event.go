package websocket

import (
	"context"
	"time"
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

// Ping 心跳检测
func (event *Events) Ping(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	var currentTime = uint64(time.Now().Unix())
	client.Heartbeat(currentTime)

	return codeStatusOk, codeMessageOk, "pong"
}
