package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "mt/api/v1/account"
	"mt/internal/app"
	"mt/internal/biz"
	typeAccount "mt/internal/constant/types/account"
	"time"
)

type AccountService struct {
	pb.UnimplementedServiceServer

	uc *biz.AccountUsecase

	tools *app.Tools
}

func NewAccountService(uc *biz.AccountUsecase, tools *app.Tools) *AccountService {
	return &AccountService{uc: uc, tools: tools}
}

// Create 创建账号
func (s *AccountService) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	createRequest := &typeAccount.CreateRequest{
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
		CreatedAt: createResponse.CreatedAt.Unix(),
	}

	return resp, nil
}

// Update 更新账号
func (s *AccountService) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	updateRequest := &typeAccount.UpdateRequest{
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
		CreatedAt: updateResponse.CreatedAt.Unix(),
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

// GetInfo 获取账号信息
func (s *AccountService) GetInfo(ctx context.Context, req *pb.GetInfoRequest) (*pb.GetInfoResponse, error) {
	getInfoResponse, err := s.uc.GetInfo(ctx, req.GetAccountId())
	if err != nil {
		return nil, err
	}

	resp := &pb.GetInfoResponse{
		AccountId:      getInfoResponse.AccountId,
		Nickname:       getInfoResponse.Nickname,
		Avatar:         getInfoResponse.Avatar,
		IsAdmin:        getInfoResponse.IsAdmin,
		IsOnline:       getInfoResponse.IsOnline,
		LastLoginIp:    getInfoResponse.LastLoginIp,
		FirstLoginTime: getInfoResponse.FirstLoginTime.Unix(),
		LastLoginTime:  getInfoResponse.LastLoginTime.Unix(),
		CreatedAt:      getInfoResponse.CreatedAt.Unix(),
		UpdatedAt:      getInfoResponse.UpdatedAt.Unix(),
		DeletedAt:      getInfoResponse.DeletedAt.Unix(),
	}

	return resp, nil
}

// UpdateLogin 更新帐号登录信息
func (s *AccountService) UpdateLogin(ctx context.Context, req *pb.UpdateLoginRequest) (*pb.UpdateLoginResponse, error) {
	updateLoginRequest := &typeAccount.UpdateLoginRequest{
		ClientIp: req.ClientIp,
	}

	updateLoginResponse, err := s.uc.UpdateLogin(ctx, req.GetAccountId(), updateLoginRequest)
	if err != nil {
		return nil, err
	}

	resp := &pb.UpdateLoginResponse{
		AccountId:      updateLoginResponse.AccountId,
		Nickname:       updateLoginResponse.Nickname,
		Avatar:         updateLoginResponse.Avatar,
		IsAdmin:        updateLoginResponse.IsAdmin,
		IsOnline:       updateLoginResponse.IsOnline,
		LastLoginIp:    updateLoginResponse.LastLoginIp,
		FirstLoginTime: updateLoginResponse.FirstLoginTime.Unix(),
		LastLoginTime:  updateLoginResponse.LastLoginTime.Unix(),
	}

	return resp, nil
}

// GenerateToken 生成TOKEN
func (s *AccountService) GenerateToken(ctx context.Context, req *pb.GenerateTokenRequest) (*pb.GenerateTokenResponse, error) {
	var nowTime = time.Now()
	generateTokenResponse, err := s.uc.GenerateToken(ctx, req.GetAccountId(), req.GetTtl())
	if err != nil {
		return nil, err
	}

	resp := &pb.GenerateTokenResponse{
		AccountId:   generateTokenResponse.AccountId,
		Token:       generateTokenResponse.Token,
		TokenExpire: nowTime.Add(time.Duration(generateTokenResponse.TokenExpire) * time.Second).Unix(),
	}

	return resp, nil
}
