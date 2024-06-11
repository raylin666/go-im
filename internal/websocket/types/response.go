package types

import "encoding/json"

type Response struct {
	Seq      string             `json:"seq"`      // 消息唯一ID, 用于标识服务端回复的是哪一条信息
	Event    string             `json:"event"`    // 消息事件
	Response *responseAgreement `json:"response"` // 消息内容
}

type responseAgreement struct {
	Code    uint32      `json:"code"`    // 响应状态码
	Message string      `json:"message"` // 响应描述
	Data    interface{} `json:"data"`    // JSON 数据包
}

// NewResponse 创建返回消息
func NewResponse(seq string, event string, code uint32, message string, data interface{}) *Response {
	response := buildResponse(code, message, data)
	return &Response{Seq: seq, Event: event, Response: response}
}

func (resp *Response) String() (headStr string) {
	headBytes, _ := json.Marshal(resp)
	headStr = string(headBytes)
	return
}

func buildResponse(code uint32, message string, data interface{}) *responseAgreement {
	return &responseAgreement{Code: code, Message: message, Data: data}
}
