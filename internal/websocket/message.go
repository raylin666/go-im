package websocket

import "encoding/json"

// Message 消息回包数据格式
type Message struct {
	Event string      // 消息事件
	Code  uint32      // 响应状态码
	Msg   string      // 响应描述
	Data  interface{} // JSON 数据包
}

// MessageRequest 消息请求数据格式
type MessageRequest struct {
	Seq   string      `json:"seq"`            // 消息唯一标识ID, 用于标识服务端回复的是哪条信息(可能有多个回复), 需客户端保证消息ID的唯一性
	Event string      `json:"event"`          // 消息事件
	Data  interface{} `json:"data,omitempty"` // JSON 数据包
}

// MessageResponse 消息响应数据格式
type MessageResponse struct {
	Seq      string                    `json:"seq"`      // 消息唯一标识ID, 用于标识服务端回复的是哪条信息(可能有多个回复)
	Event    string                    `json:"event"`    // 消息事件
	Response *MessageResponseAgreement `json:"response"` // 消息内容
}

type MessageResponseAgreement struct {
	Code    uint32      `json:"code"`    // 响应状态码
	Message string      `json:"message"` // 响应描述
	Data    interface{} `json:"data"`    // JSON 数据包
}

// NewMessageResponse 创建响应消息数据包
func NewMessageResponse(seq string, event string, code uint32, message string, data interface{}) *MessageResponse {
	var response = &MessageResponseAgreement{
		Code:    code,
		Message: message,
		Data:    data,
	}

	return &MessageResponse{Seq: seq, Event: event, Response: response}
}

// String 消息转换为字符串
func (message *MessageResponse) String() (headStr string) {
	headBytes, _ := json.Marshal(message)
	headStr = string(headBytes)
	return
}
