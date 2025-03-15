package websocket

import (
	"context"
	"encoding/json"
	"net/http"
)

type C2CMessageRequest struct {
	// 接收者账号ID
	ToAccount string `json:"to"`
	// 消息内容
	Message interface{} `json:"message"`
}

// C2CMessage 发送C2C消息
func (event *messageEvent) C2CMessage(ctx context.Context, client *Client, seq string, message []byte) (messages []Message) {
	//TODO implement me

	// TODO 数据包合法性校验/解析消息数据包
	request := &C2CMessageRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {
		code := http.StatusUnprocessableEntity
		messages = append(messages, Message{Event: MessageEventC2CMessage, Code: uint32(code), Msg: http.StatusText(code), Data: "C2C消息事件数据包协议格式错误"})
		return
	}

	if request.ToAccount == "" {

	}

	// TODO 处理发送给接收者消息

	return
}
