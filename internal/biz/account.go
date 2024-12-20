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
	Login(ctx context.Context, accountId string, data *typeAccount.LoginRequest) (*model.Account, *model.AccountOnline, error)
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
	account, err := uc.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	resp := &typeAccount.CreateResponse{
		AccountId: account.AccountId,
		Nickname:  account.Nickname,
		Avatar:    account.Avatar,
		IsAdmin:   account.IsAdmin == 1,
		CreatedAt: account.CreatedAt,
	}

	return resp, nil
}

// Update 更新账号
func (uc *AccountUsecase) Update(ctx context.Context, accountId string, req *typeAccount.UpdateRequest) (*typeAccount.UpdateResponse, error) {
	account, err := uc.repo.Update(ctx, accountId, req)
	if err != nil {
		return nil, err
	}

	resp := &typeAccount.UpdateResponse{
		AccountId: account.AccountId,
		Nickname:  account.Nickname,
		Avatar:    account.Avatar,
		IsAdmin:   account.IsAdmin == 1,
		CreatedAt: account.CreatedAt,
	}

	return resp, nil
}

// Delete 删除账号
func (uc *AccountUsecase) Delete(ctx context.Context, accountId string) (*typeAccount.DeleteResponse, error) {
	account, err := uc.repo.Delete(ctx, accountId)
	if err != nil {
		return nil, err
	}

	return &typeAccount.DeleteResponse{AccountId: account.AccountId}, nil
}

// GetInfo 获取账号信息
func (uc *AccountUsecase) GetInfo(ctx context.Context, accountId string) (*typeAccount.GetInfoResponse, error) {
	account, err := uc.repo.GetInfo(ctx, accountId)
	if err != nil {
		return nil, err
	}

	resp := &typeAccount.GetInfoResponse{
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
func (uc *AccountUsecase) Login(ctx context.Context, accountId string, req *typeAccount.LoginRequest) (*typeAccount.LoginResponse, error) {
	account, accountOnline, err := uc.repo.Login(ctx, accountId, req)
	if err != nil {
		return nil, err
	}

	resp := &typeAccount.LoginResponse{
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

// GenerateToken 生成TOKEN
func (uc *AccountUsecase) GenerateToken(ctx context.Context, accountId string, ttl int64) (*typeAccount.GenerateTokenResponse, error) {
	if accountId == "" {
		return nil, defined.ErrorRequestParams
	}

	// 默认Token为1天过期
	if ttl <= 0 {
		ttl = 86400
	}

	token, err := uc.tools.JWT().GenerateToken(accountId, time.Duration(ttl)*time.Second, auth.JWTClaimsOptions{})
	if err != nil {
		return nil, defined.ErrorGenerateToken
	}

	resp := &typeAccount.GenerateTokenResponse{
		AccountId:   accountId,
		Token:       token,
		TokenExpire: ttl,
	}

	return resp, nil
}
