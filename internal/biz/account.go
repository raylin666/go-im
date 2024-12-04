package biz

import (
	"context"
	"github.com/raylin666/go-utils/auth"
	"mt/internal/app"
	"mt/internal/constant/defined"
	typeAccount "mt/internal/constant/types/account"
	"mt/internal/repositories/dbrepo/model"
	"time"
)

type Account struct {
}

type AccountRepo interface {
	Create(ctx context.Context, data *typeAccount.CreateRequest) (*model.Account, error)
	Update(ctx context.Context, accountId string, data *typeAccount.UpdateRequest) (*model.Account, error)
	Delete(ctx context.Context, accountId string) (*model.Account, error)
	GetInfo(ctx context.Context, accountId string) (*model.Account, error)
	UpdateLogin(ctx context.Context, accountId string, data *typeAccount.UpdateLoginRequest) (*model.Account, error)
}

type AccountUsecase struct {
	repo  AccountRepo
	tools *app.Tools
}

func NewAccountUsecase(repo AccountRepo, tools *app.Tools) *AccountUsecase {
	return &AccountUsecase{repo: repo, tools: tools}
}

// Create 创建账号
func (uc *AccountUsecase) Create(ctx context.Context, req *typeAccount.CreateRequest) (*typeAccount.CreateResponse, error) {
	m, err := uc.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := &typeAccount.CreateResponse{
		AccountId: m.AccountId,
		Nickname:  m.Nickname,
		Avatar:    m.Avatar,
		IsAdmin:   m.IsAdmin == 1,
		CreatedAt: m.CreatedAt,
	}

	return resp, nil
}

// Update 更新账号
func (uc *AccountUsecase) Update(ctx context.Context, accountId string, req *typeAccount.UpdateRequest) (*typeAccount.UpdateResponse, error) {
	m, err := uc.repo.Update(ctx, accountId, req)
	if err != nil {
		return nil, err
	}

	resp := &typeAccount.UpdateResponse{
		AccountId: m.AccountId,
		Nickname:  m.Nickname,
		Avatar:    m.Avatar,
		IsAdmin:   m.IsAdmin == 1,
		CreatedAt: m.CreatedAt,
	}

	return resp, nil
}

// Delete 删除账号
func (uc *AccountUsecase) Delete(ctx context.Context, accountId string) (*typeAccount.DeleteResponse, error) {
	m, err := uc.repo.Delete(ctx, accountId)
	if err != nil {
		return nil, err
	}

	return &typeAccount.DeleteResponse{AccountId: m.AccountId}, nil
}

// GetInfo 获取账号信息
func (uc *AccountUsecase) GetInfo(ctx context.Context, accountId string) (*typeAccount.GetInfoResponse, error) {
	m, err := uc.repo.GetInfo(ctx, accountId)
	if err != nil {
		return nil, err
	}

	resp := &typeAccount.GetInfoResponse{
		AccountId:      m.AccountId,
		Nickname:       m.Nickname,
		Avatar:         m.Avatar,
		IsAdmin:        m.IsAdmin == 1,
		IsOnline:       m.IsOnline == 1,
		LastLoginIp:    m.LastLoginIp,
		FirstLoginTime: m.FirstLoginTime,
		LastLoginTime:  m.LastLoginTime,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      &m.DeletedAt.Time,
	}

	return resp, nil
}

// UpdateLogin 更新帐号登录信息
func (uc *AccountUsecase) UpdateLogin(ctx context.Context, accountId string, req *typeAccount.UpdateLoginRequest) (*typeAccount.UpdateLoginResponse, error) {
	m, err := uc.repo.UpdateLogin(ctx, accountId, req)
	if err != nil {
		return nil, err
	}

	resp := &typeAccount.UpdateLoginResponse{
		AccountId:      m.AccountId,
		Nickname:       m.Nickname,
		Avatar:         m.Avatar,
		IsAdmin:        m.IsAdmin == 1,
		IsOnline:       m.IsOnline == 1,
		LastLoginIp:    m.LastLoginIp,
		FirstLoginTime: m.FirstLoginTime,
		LastLoginTime:  m.LastLoginTime,
	}

	return resp, nil
}

// GenerateToken 生成TOKEN
func (uc *AccountUsecase) GenerateToken(ctx context.Context, accountId string, ttl int64) (*typeAccount.GenerateTokenResponse, error) {
	if accountId == "" {
		return nil, defined.ErrorRequestParamsError
	}

	// 默认Token为1天过期
	if ttl <= 0 {
		ttl = 86400
	}

	token, err := uc.tools.JWT().GenerateToken(accountId, time.Duration(ttl)*time.Second, auth.JWTClaimsOptions{})
	if err != nil {
		return nil, defined.ErrorGenerateTokenError
	}

	resp := &typeAccount.GenerateTokenResponse{
		AccountId:   accountId,
		Token:       token,
		TokenExpire: ttl,
	}

	return resp, nil
}
