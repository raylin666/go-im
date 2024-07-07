package event

import "mt/internal/websocket/types/message"

// SendC2CMessageRequest 事件C2C消息发送请求
type SendC2CMessageRequest struct {
	// 消息接收者ID
	ToAccount string `json:"to_account"`

	/**
	消息类型
		文本消息: MessageTypeText
	*/
	Type string `json:"type"`

	// 消息内容
	Content string `json:"content"`
}

// SendC2CMessageResponse 事件C2C消息发送响应
type SendC2CMessageResponse struct {
	// 消息ID
	MessageId int `json:"message_id"`
	// 消息发送者ID
	FromAccount string `json:"from_account"`
	// 消息接收者ID
	ToAccount string `json:"to_account"`
	// 消息发送时间戳 单位/秒
	MessageSendTime int64 `json:"message_send_time"`

	/**
	消息类型
		文本消息: MessageTypeText
	*/
	Type string `json:"type"`

	// 消息内容 JSON 数据包
	Content message.BuilderType `json:"content"`
}
