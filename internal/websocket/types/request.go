package types

// Request 请求数据格式
type Request struct {
	Seq   string      `json:"seq"`            // 消息唯一ID, 用于标识服务端回复的是哪一条信息
	Event string      `json:"event"`          // 消息事件
	Data  interface{} `json:"data,omitempty"` // JSON 数据包
}
