package websocket

import (
	"context"
)

// Bind 客户端和账号信息绑定
func (event *messageEvent) Bind(ctx context.Context, client *Client, seq string, message []byte) []Message {
	//TODO implement me
	panic("implement me")
}
