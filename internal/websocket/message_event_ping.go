package websocket

import (
	"context"
	"time"
)

// Ping 心跳检测
func (event *messageEvent) Ping(ctx context.Context, client *Client, seq string, message []byte) (messages []Message) {
	//TODO implement me

	// 更新客户端连接心跳
	client.Heartbeat(time.Now())

	// 发送 PONE 回包
	messages = append(messages, Message{Event: MessageEventPing, Data: "PONE"})

	return messages
}
