package event

import (
	"context"
	"encoding/json"
	"mt/internal/constant/defined"
	"mt/internal/websocket"
	"mt/internal/websocket/types"
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

	data = &client.Account

	js, _ := json.Marshal(&types.Request{Seq: "1", Event: "ping"})
	client.EventMessageHandler(ctx, js)

	return
}
