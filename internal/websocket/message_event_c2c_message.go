package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	kratosErr "github.com/go-kratos/kratos/v2/errors"
	"mt/internal/constant/types"
	"net/http"
)

type C2CMessageRequest struct {
	// 接收者账号ID
	ToAccount string `json:"to"`
	// 消息内容
	Message string `json:"message"`
}

type C2CMessageResponse struct {
	// 发送者账号ID
	FromAccount string `json:"from"`
	// 发送者昵称
	FromNickname string `json:"from_nickname"`
	// 发送者头像
	FromAvatar string `json:"from_avatar"`
	// 接收者账号ID
	ToAccount string `json:"to"`
	// 接收者昵称
	ToNickname string `json:"to_nickname"`
	// 接收者头像
	ToAvatar string `json:"to_avatar"`
	// 消息内容
	Message string `json:"message"`
}

// C2CMessage 发送C2C消息
func (event *messageEvent) C2CMessage(ctx context.Context, client *Client, seq string, message []byte) (messages []Message) {
	//TODO implement me

	var errTitle = "C2C消息事件"

	// TODO 数据包合法性校验/解析消息数据包
	request := &C2CMessageRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
		code := http.StatusUnprocessableEntity
		messages = append(messages, Message{Event: MessageEventC2CMessage, Code: uint32(code), Msg: http.StatusText(code), Data: fmt.Sprintf("%s数据包协议格式错误", errTitle)})
		return
	}

	// TODO 调用发送 C2C 消息
	dataLogicResult := event.DataLogicRepo.Message.SendC2CMessage(ctx, &types.MessageSendC2CMessageRequest{
		Seq:         seq,
		FromAccount: client.Account.ID,
		ToAccount:   request.ToAccount,
		Message:     request.Message,
	})

	if err != nil {
		errDetail := kratosErr.FromError(err)
		code := errDetail.Code
		errData := fmt.Sprintf("%s处理错误", errTitle)
		if errDetail.Message != "" {
			errData = errData + " - " + errDetail.Message
		}

		messages = append(messages, Message{Event: MessageEventC2CMessage, Code: uint32(code), Msg: http.StatusText(int(code)), Data: errData})
		return
	}

	messageData := C2CMessageResponse{
		FromAccount:  dataLogicResult.FromAccount.AccountId,
		FromNickname: dataLogicResult.FromAccount.Nickname,
		FromAvatar:   dataLogicResult.FromAccount.Avatar,
		ToAccount:    dataLogicResult.ToAccount.AccountId,
		ToNickname:   dataLogicResult.ToAccount.Nickname,
		ToAvatar:     dataLogicResult.ToAccount.Avatar,
		Message:      dataLogicResult.Message,
	}

	// TODO 消息回包给发送者
	messages = append(messages, Message{Event: MessageEventC2CMessage, Data: messageData})

	// TODO 消息发送给接收者
	if dataLogicResult.ToAccountOnline {
		
	}

	return
}
