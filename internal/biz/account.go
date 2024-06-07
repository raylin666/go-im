package biz

import (
	"context"
	"mt/internal/constant/types"
	"mt/internal/repositories/dbrepo/model"
	"mt/pkg/logger"
)

type Account struct {
}

type AccountRepo interface {
	Create(ctx context.Context, data types.AccountCreateData) (*model.Account, error)
	Update(ctx context.Context, data types.AccountUpdateData) (*model.Account, error)
	Delete(ctx context.Context, accountId string) (*model.Account, error)
}

type AccountUsecase struct {
	repo AccountRepo
	log  *logger.Logger
}

func NewAccountUsecase(repo AccountRepo, logger *logger.Logger) *AccountUsecase {
	return &AccountUsecase{repo: repo, log: logger}
}

// Create 创建账号
func (uc *AccountUsecase) Create(ctx context.Context, req *types.AccountCreateRequest) (*types.AccountCreateResponse, error) {
	var createData types.AccountCreateData
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

	resp := &types.AccountCreateResponse{
		AccountId: m.AccountId,
		Nickname:  m.Nickname,
		Avatar:    m.Avatar,
		IsAdmin:   m.IsAdmin == 1,
		CreatedAt: m.CreatedAt,
	}

	return resp, nil
}

// Update 更新账号
func (uc *AccountUsecase) Update(ctx context.Context, accountId string, req *types.AccountUpdateRequest) (*types.AccountUpdateResponse, error) {
	var updateData types.AccountUpdateData
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

	resp := &types.AccountUpdateResponse{
		AccountId: m.AccountId,
		Nickname:  m.Nickname,
		Avatar:    m.Avatar,
		IsAdmin:   m.IsAdmin == 1,
		CreatedAt: m.CreatedAt,
	}

	return resp, nil
}

// Delete 删除账号
func (uc *AccountUsecase) Delete(ctx context.Context, accountId string) (*types.AccountDeleteResponse, error) {
	m, err := uc.repo.Delete(ctx, accountId)
	if err != nil {
		return nil, err
	}

	return &types.AccountDeleteResponse{AccountId: m.AccountId}, nil
}
