package websocket

import (
	"context"
)

const (
	EventPing        = "ping"
	EventBind        = "bind"
	EventAccountInfo = "account_info"
)

// Events 所有消息事件接口
type Events interface {
	// GetAll 获取所有事件对应处理器
	GetAll() map[string]EventDisposeFunc
	// ClientWhiteEventNames 客户端请求所支持的事件, 不在指定的事件客户端无法调用
	ClientWhiteEventNames() []string

	// Ping 心跳检测
	Ping(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}, send bool)
	// Bind 客户端和账号信息绑定
	Bind(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}, send bool)
	// AccountInfo 获取账号信息
	AccountInfo(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}, send bool)
}

type EventDisposeFunc func(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}, send bool)
