package event

import (
	"context"
	"mt/internal/websocket"
	"time"
)

var _ websocket.MessageEvent = (*Event)(nil)

type Event struct{}

// Ping 心跳检测[消息事件处理]
func (event Event) Ping(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code, msg, _ = defaultEventResponse()
	data = "pong"

	client.Heartbeat(time.Now())

	return
}

// defaultEventResponse 默认事件返回值
func defaultEventResponse() (code uint32, msg string, data interface{}) {
	code, msg, data = 200, "OK", nil
	return
}
