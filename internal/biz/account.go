package biz

import (
	"context"
	"github.com/raylin666/go-utils/auth"
	"mt/errors"
	"mt/internal/app"
	"mt/internal/constant/types"
	"mt/internal/data"
	"time"
)

type AccountUsecase struct {
	repo  data.AccountRepo
	tools *app.Tools
}

func NewAccountUsecase(repo data.AccountRepo, tools *app.Tools) *AccountUsecase {
	return &AccountUsecase{repo: repo, tools: tools}
}

// Create 创建账号
func (uc *AccountUsecase) Create(ctx context.Context, req *types.AccountCreateRequest) (*types.AccountCreateResponse, error) {
	account, err := uc.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := &types.AccountCreateResponse{
		AccountId: account.AccountId,
		Nickname:  account.Nickname,
		Avatar:    account.Avatar,
		IsAdmin:   account.IsAdmin == 1,
		CreatedAt: account.CreatedAt,
	}

	return resp, nil
}

// Update 更新账号
func (uc *AccountUsecase) Update(ctx context.Context, accountId string, req *types.AccountUpdateRequest) (*types.AccountUpdateResponse, error) {
	account, err := uc.repo.Update(ctx, accountId, req)
	if err != nil {
		return nil, err
	}

	resp := &types.AccountUpdateResponse{
		AccountId: account.AccountId,
		Nickname:  account.Nickname,
		Avatar:    account.Avatar,
		IsAdmin:   account.IsAdmin == 1,
		CreatedAt: account.CreatedAt,
	}

	return resp, nil
}

// Delete 删除账号
func (uc *AccountUsecase) Delete(ctx context.Context, accountId string) (*types.AccountDeleteResponse, error) {
	account, err := uc.repo.Delete(ctx, accountId)
	if err != nil {
		return nil, err
	}

	return &types.AccountDeleteResponse{AccountId: account.AccountId}, nil
}

// GetInfo 获取账号信息
func (uc *AccountUsecase) GetInfo(ctx context.Context, accountId string) (*types.AccountGetInfoResponse, error) {
	account, err := uc.repo.GetInfo(ctx, accountId)
	if err != nil {
		return nil, err
	}

	resp := &types.AccountGetInfoResponse{
		AccountId:      account.AccountId,
		Nickname:       account.Nickname,
		Avatar:         account.Avatar,
		IsAdmin:        account.IsAdmin == 1,
		IsOnline:       account.IsOnline == 1,
		LastLoginIp:    account.LastLoginIp,
		FirstLoginTime: account.FirstLoginTime,
		LastLoginTime:  account.LastLoginTime,
		CreatedAt:      account.CreatedAt,
		UpdatedAt:      account.UpdatedAt,
		DeletedAt:      &account.DeletedAt.Time,
	}

	return resp, nil
}

// Login 登录帐号
func (uc *AccountUsecase) Login(ctx context.Context, accountId string, req *types.AccountLoginRequest) (*types.AccountLoginResponse, error) {
	account, accountOnline, err := uc.repo.Login(ctx, accountId, req)
	if err != nil {
		return nil, err
	}

	resp := &types.AccountLoginResponse{
		AccountId:      account.AccountId,
		Nickname:       account.Nickname,
		Avatar:         account.Avatar,
		IsAdmin:        account.IsAdmin == 1,
		IsOnline:       account.IsOnline == 1,
		LastLoginIp:    account.LastLoginIp,
		FirstLoginTime: account.FirstLoginTime,
		LastLoginTime:  account.LastLoginTime,
		OnlineId:       accountOnline.ID,
	}

	return resp, nil
}

// Logout 登出帐号
func (uc *AccountUsecase) Logout(ctx context.Context, accountId string, req *types.AccountLogoutRequest) error {
	_, err := uc.repo.Logout(ctx, accountId, req)
	return err
}

// GenerateToken 生成TOKEN
func (uc *AccountUsecase) GenerateToken(ctx context.Context, accountId string, ttl int64) (*types.AccountGenerateTokenResponse, error) {
	if accountId == "" {
		return nil, errors.New().RequestParams()
	}

	// 默认Token为1天过期
	if ttl <= 0 {
		ttl = 86400
	}

	token, err := uc.tools.JWT().GenerateToken(accountId, time.Duration(ttl)*time.Second, auth.JWTClaimsOptions{})
	if err != nil {
		return nil, errors.New().GenerateToken()
	}

	resp := &types.AccountGenerateTokenResponse{
		AccountId:   accountId,
		Token:       token,
		TokenExpire: ttl,
	}

	return resp, nil
}
