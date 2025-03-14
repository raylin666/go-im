package websocket

import (
	"context"
)

// C2CMessage 发送C2C消息
func (event *messageEvent) C2CMessage(ctx context.Context, client *Client, seq string, message []byte) []Message {
	//TODO implement me
	panic("implement me")
}
