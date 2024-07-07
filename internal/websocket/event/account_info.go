package event

import (
	"context"
	"mt/internal/websocket"
	typesEvent "mt/internal/websocket/types/event"
)

// AccountInfo 获取账号信息 [消息事件处理]
func (event *events) AccountInfo(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}, send bool) {
	code, msg, _, send = defaultEventResponse()

	data = typesEvent.AccountInfoResponse{
		ID:       client.Account.ID,
		Nickname: client.Account.Nickname,
		Avatar:   client.Account.Avatar,
	}

	return
}
