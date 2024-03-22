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

func (uc *AccountUsecase) Create(ctx context.Context, req *types.AccountCreateRequest) (*types.AccountCreateResponse, error) {
	var createData types.AccountCreateData
	createData.UserId = req.UserId
	createData.Username = req.Username
	createData.Avatar = req.Avatar

	if req.IsAdmin {
		createData.IsAdmin = 1
	} else {
		createData.IsAdmin = 0
	}

	m, err := uc.repo.Create(ctx, createData)
	if err != nil {
		return nil, err
	}

	resp := &types.AccountCreateResponse{
		UserId:    m.UserId,
		Username:  m.Username,
		Avatar:    m.Avatar,
		IsAdmin:   m.IsAdmin == 1,
		CreatedAt: m.CreatedAt,
	}

	return resp, nil
}
