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
	Create(context.Context, types.AccountCreateData) (*model.Account, error)
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
