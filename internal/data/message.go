package data

import (
	"context"
	"fmt"
	"mt/errors"
	"mt/internal/app"
	"mt/internal/constant/types"
	"mt/internal/repositories"
)

type MessageRepo interface {
	SendC2CMessage(ctx context.Context, data *types.MessageSendC2CMessageRequest) error
}

type messageRepo struct {
	data  repositories.DataRepo
	tools *app.Tools
}

func NewMessageRepo(repo repositories.DataRepo, tools *app.Tools) MessageRepo {
	return &messageRepo{
		data:  repo,
		tools: tools,
	}
}

// SendC2CMessage 发送 C2C 消息
func (r *messageRepo) SendC2CMessage(ctx context.Context, data *types.MessageSendC2CMessageRequest) error {
	if data.Message == "" {
		return errors.New().SendMessageContentRequired()
	}

	if data.ToAccount == "" {
		return errors.New().ToAccountNotFound()
	}

	fmt.Println(data)

	return nil
}
