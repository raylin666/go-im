package service

import (
	"context"
	pb "mt/api/v1/account"
	"mt/internal/biz"
	"mt/internal/constant/types"
)

type AccountService struct {
	pb.UnimplementedAccountServer

	uc *biz.AccountUsecase
}

func NewAccountService(uc *biz.AccountUsecase) *AccountService {
	return &AccountService{uc: uc}
}

// Create 创建账号
func (s *AccountService) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	// 权限验证

	createRequest := &types.AccountCreateRequest{
		UserId:   req.GetUserId(),
		Username: req.GetUsername(),
		Avatar:   req.GetAvatar(),
		Status:   int8(req.GetStatus()),
		IsAdmin:  req.GetIsAdmin(),
	}

	_, err := s.uc.Create(ctx, createRequest)
	if err != nil {
		return nil, err
	}

	resp := &pb.CreateResponse{}

	return resp, nil
}
