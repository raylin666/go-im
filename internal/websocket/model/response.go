package model

import "encoding/json"

type Head struct {
	Seq      string    `json:"seq"`      // 消息唯一ID
	Event    string    `json:"event"`    // 消息事件
	Response *Response `json:"response"` // 消息内容
}

type Response struct {
	Code    uint32      `json:"code"`    // 响应状态码
	Message string      `json:"message"` // 响应描述
	Data    interface{} `json:"data"`    // JSON 数据包
}

// NewResponseHead 创建返回消息
func NewResponseHead(seq string, event string, code uint32, message string, data interface{}) *Head {
	response := NewResponse(code, message, data)
	return &Head{Seq: seq, Event: event, Response: response}
}

func (h *Head) String() (headStr string) {
	headBytes, _ := json.Marshal(h)
	headStr = string(headBytes)
	return
}

func NewResponse(code uint32, message string, data interface{}) *Response {
	return &Response{Code: code, Message: message, Data: data}
}
