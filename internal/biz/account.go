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
	Create(ctx context.Context, data typeAccount.CreateData) (*model.Account, error)
	Update(ctx context.Context, data typeAccount.UpdateData) (*model.Account, error)
	Delete(ctx context.Context, accountId string) (*model.Account, error)
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
	var createData typeAccount.CreateData
	createData.AccountId = req.AccountId
	createData.Nickname = req.Nickname
	createData.Avatar = req.Avatar

	createData.IsAdmin = 0
	if req.IsAdmin {
		createData.IsAdmin = 1
	}

	m, err := uc.repo.Create(ctx, createData)
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
	var updateData typeAccount.UpdateData
	updateData.AccountId = accountId
	updateData.Nickname = req.Nickname
	updateData.Avatar = req.Avatar

	updateData.IsAdmin = 0
	if req.IsAdmin {
		updateData.IsAdmin = 1
	}

	m, err := uc.repo.Update(ctx, updateData)
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
