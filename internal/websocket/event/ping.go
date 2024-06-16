package event

import (
	"context"
	"mt/internal/websocket"
	"time"
)

// Ping 心跳检测[消息事件处理]
func (event *events) Ping(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}, send bool) {
	code, msg, _, send = defaultEventResponse()
	data = "pong"

	client.Heartbeat(time.Now())

	return
}
