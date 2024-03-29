package service

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "mt/api/v1/account"
	"mt/internal/biz"
	"mt/internal/constant/types"
)

type AccountService struct {
	pb.UnimplementedServiceServer

	uc *biz.AccountUsecase
}

func NewAccountService(uc *biz.AccountUsecase) *AccountService {
	return &AccountService{uc: uc}
}

// Create 创建账号
func (s *AccountService) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	createRequest := &types.AccountCreateRequest{
		UserId:   req.GetUserId(),
		Username: req.GetUsername(),
		Avatar:   req.GetAvatar(),
		IsAdmin:  req.GetIsAdmin(),
	}

	createResponse, err := s.uc.Create(ctx, createRequest)
	if err != nil {
		return nil, err
	}

	resp := &pb.CreateResponse{
		UserId:    createResponse.UserId,
		Username:  createResponse.Username,
		Avatar:    createResponse.Avatar,
		IsAdmin:   createResponse.IsAdmin,
		CreatedAt: timestamppb.New(createResponse.CreatedAt),
	}

	return resp, nil
}
