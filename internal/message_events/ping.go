package websocket_events

import (
	"context"
	"mt/internal/websocket"
)

// Ping 心跳检测
func (event *messageEvent) Ping(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	//TODO implement me
	panic("implement me")
}
