package service

import (
	"context"
	pb "mt/api/v1/message"
	"mt/internal/app"
	"mt/internal/biz"
	"mt/internal/constant/types"
)

type MessageService struct {
	pb.UnimplementedServiceServer

	uc *biz.MessageUsecase

	tools *app.Tools
}

func NewMessageService(uc *biz.MessageUsecase, tools *app.Tools) *MessageService {
	return &MessageService{uc: uc, tools: tools}
}

// SendC2CMessage 发送 C2C 消息
func (s *MessageService) SendC2CMessage(ctx context.Context, req *pb.SendC2CMessageRequest) (*pb.SendC2CMessageResponse, error) {
	sendC2CMessageRequest := &types.MessageSendC2CMessageRequest{
		Seq:       req.GetSeq(),
		ToAccount: req.GetToAccount(),
		Message:   req.GetMessage(),
	}

	_, err := s.uc.SendC2CMessage(ctx, sendC2CMessageRequest)
	if err != nil {
		return nil, err
	}

	resp := &pb.SendC2CMessageResponse{}

	return resp, nil
}
