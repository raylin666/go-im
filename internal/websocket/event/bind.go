package event

import (
	"context"
	"mt/internal/websocket"
)

// Bind 客户端和账号信息绑定
func (event *messageEvent) Bind(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	//TODO implement me
	panic("implement me")
}
