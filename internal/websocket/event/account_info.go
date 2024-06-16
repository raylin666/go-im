package event

import (
	"context"
	"mt/internal/websocket"
)

// AccountInfo 获取账号信息 [消息事件处理]
func (event *events) AccountInfo(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}, send bool) {
	code, msg, _, send = defaultEventResponse()

	data = client.Account

	return
}
