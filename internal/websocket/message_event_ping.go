package websocket

import (
	"context"
)

// Ping 心跳检测
func (event *messageEvent) Ping(ctx context.Context, client *Client, seq string, message []byte) (messages []Message) {
	//TODO implement me

	// 发送 PONE 回包
	messages = append(messages, Message{Event: MessageEventPing, Data: "PONE"})

	return
}
