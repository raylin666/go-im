package event

import (
	"context"
	"mt/internal/constant/defined"
	"mt/internal/websocket"
	"mt/pkg/utils"
)

// Bind 客户端和账号信息绑定 [消息事件处理]
func (event *events) Bind(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}, send bool) {
	code, msg, _, send = defaultEventResponse()
	send = false // Bind 消息事件不需要推送给客户端

	q := event.dbQuery().Account
	account, err := q.WithContext(ctx).FirstByAccountId(client.Account.ID)
	if err != nil {
		code, msg = utils.ErrorMessage(defined.ErrorNotVisitAuth)
		return
	}

	// 更新账号信息
	client.Account.ID = account.AccountId
	client.Account.Nickname = account.Nickname
	client.Account.Avatar = account.Avatar

	// 推送账号信息消息事件
	event.NewPushMessage(ctx, client, websocket.EventAccountInfo, seq, nil)

	return
}
