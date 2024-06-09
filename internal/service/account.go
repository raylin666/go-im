package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
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
		AccountId: req.GetAccountId(),
		Nickname:  req.GetNickname(),
		Avatar:    req.GetAvatar(),
		IsAdmin:   req.GetIsAdmin(),
	}

	createResponse, err := s.uc.Create(ctx, createRequest)
	if err != nil {
		return nil, err
	}

	resp := &pb.CreateResponse{
		AccountId: createResponse.AccountId,
		Nickname:  createResponse.Nickname,
		Avatar:    createResponse.Avatar,
		IsAdmin:   createResponse.IsAdmin,
		CreatedAt: timestamppb.New(createResponse.CreatedAt),
	}

	return resp, nil
}

// Update 更新账号
func (s *AccountService) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	updateRequest := &types.AccountUpdateRequest{
		Nickname: req.GetNickname(),
		Avatar:   req.GetAvatar(),
		IsAdmin:  req.GetIsAdmin(),
	}

	updateResponse, err := s.uc.Update(ctx, req.GetAccountId(), updateRequest)
	if err != nil {
		return nil, err
	}

	resp := &pb.UpdateResponse{
		AccountId: updateResponse.AccountId,
		Nickname:  updateResponse.Nickname,
		Avatar:    updateResponse.Avatar,
		IsAdmin:   updateResponse.IsAdmin,
		CreatedAt: timestamppb.New(updateResponse.CreatedAt),
	}

	return resp, nil
}

// Delete 删除账号
func (s *AccountService) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	_, err := s.uc.Delete(ctx, req.GetAccountId())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GenerateToken 生成TOKEN
func (s *AccountService) GenerateToken(ctx context.Context, req *pb.GenerateTokenRequest) (*pb.GenerateTokenResponse, error) {
	generateTokenResponse, err := s.uc.GenerateToken(ctx, req.GetAccountId(), req.GetTtl())
	if err != nil {
		return nil, err
	}

	resp := &pb.GenerateTokenResponse{
		AccountId:   generateTokenResponse.AccountId,
		Token:       generateTokenResponse.Token,
		TokenExpire: generateTokenResponse.TokenExpire,
	}

	return resp, nil
}
