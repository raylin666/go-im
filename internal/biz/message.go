package biz

import (
	"context"
	"mt/internal/app"
	"mt/internal/constant/types"
	"mt/internal/data"
	"mt/internal/websocket"
)

type MessageUsecase struct {
	repo            data.MessageRepo
	wsClientManager websocket.WSClientManager
	tools           *app.Tools
}

func NewMessageUsecase(repo data.MessageRepo, wsClientManager websocket.WSClientManager, tools *app.Tools) *MessageUsecase {
	return &MessageUsecase{repo: repo, wsClientManager: wsClientManager, tools: tools}
}

// SendC2CMessage 发送 C2C 消息
func (uc *MessageUsecase) SendC2CMessage(ctx context.Context, req *types.MessageSendC2CMessageRequest) (*types.MessageSendC2CMessageResponse, error) {
	err := uc.repo.SendC2CMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := &types.MessageSendC2CMessageResponse{}

	return resp, nil
}
