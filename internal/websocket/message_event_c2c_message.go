package websocket

import (
	"context"
)

// C2CMessage 发送C2C消息
func (event *messageEvent) C2CMessage(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	//TODO implement me
	panic("implement me")
}
