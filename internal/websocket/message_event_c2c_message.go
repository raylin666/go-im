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
	_, _, err = event.DataLogicRepo.Message.SendC2CMessage(ctx, &types.MessageSendC2CMessageRequest{
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

	return
}
