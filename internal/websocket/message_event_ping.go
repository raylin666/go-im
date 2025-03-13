package websocket

import (
	"context"
)

// Ping 心跳检测
func (event *messageEvent) Ping(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	//TODO implement me
	panic("implement me")
}
