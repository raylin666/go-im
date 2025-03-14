package websocket

import (
	"context"
	"encoding/json"
	"fmt"
)

type C2CMessageRequest struct {
	ToAccount string      `json:"to_account"`
	Message   interface{} `json:"message"`
}

// C2CMessage 发送C2C消息
func (event *messageEvent) C2CMessage(ctx context.Context, client *Client, seq string, message []byte) []Message {
	//TODO implement me

	// TODO 数据包合法性校验/解析消息数据包
	request := &C2CMessageRequest{}
	err := json.Unmarshal(message, request)
	if err != nil {

	}

	fmt.Println(string(message))

	return nil
}
