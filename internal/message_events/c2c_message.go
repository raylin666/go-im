package websocket_events

import (
	"context"
	"mt/internal/websocket"
)

// C2CMessage 发送C2C消息
func (event *messageEvent) C2CMessage(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	//TODO implement me
	panic("implement me")
}
