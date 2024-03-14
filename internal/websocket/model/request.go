package model

// Request 请求数据格式
type Request struct {
	Seq   string      `json:"seq"`            // 消息唯一ID
	Event string      `json:"event"`          // 消息事件
	Data  interface{} `json:"data,omitempty"` // JSON 数据包
}
