package websocket

import (
	"context"
	"time"
)

const (
	EventPing = "ping"
)

// Events 获取所有消息事件
func Events() (events map[string]EventDisposeFunc) {
	var event = &Event{}
	events = make(map[string]EventDisposeFunc)

	// 心跳检测
	events[EventPing] = event.Ping

	return
}

type Event struct {
}

// Ping 心跳检测[消息事件处理]
func (event *Event) Ping(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code, msg, _ = defaultEventResponse()
	data = "pong"

	client.Heartbeat(time.Now())

	return
}

type EventDisposeFunc func(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

// defaultEventResponse 默认事件返回值
func defaultEventResponse() (code uint32, msg string, data interface{}) {
	code, msg, data = 200, "OK", nil
	return
}
