package biz

import (
	"context"
	"github.com/google/uuid"
	"mt/internal/constant/types"
	"mt/internal/repositories/dbrepo/model"
	"mt/pkg/logger"
	"strings"
)

type Manager struct {
}

type ManagerRepo interface {
	Create(context.Context, types.CreateData) (*model.App, error)
}

type ManagerUsecase struct {
	repo ManagerRepo
	log  *logger.Logger
}

func NewManagerUsecase(repo ManagerRepo, logger *logger.Logger) *ManagerUsecase {
	return &ManagerUsecase{repo: repo, log: logger}
}

func (uc *ManagerUsecase) Create(ctx context.Context, req *types.CreateRequest) (*types.CreateResponse, error) {
	// 生成应用ID
	uuidApp := uuid.New()

	var createData types.CreateData
	createData.Ident = req.Ident
	createData.Name = req.Name
	createData.Key = uuidApp.ID()
	createData.Secret = strings.ReplaceAll(uuidApp.String(), "-", "")
	createData.Status = int8(req.Status)
	createData.ExpiredAt = req.ExpiredAt

	m, err := uc.repo.Create(ctx, createData)
	if err != nil {
		return nil, err
	}

	resp := &types.CreateResponse{
		Id:        m.ID,
		Ident:     m.Ident,
		Name:      m.Name,
		Key:       m.Key,
		Secret:    m.Secret,
		Status:    m.Status,
		ExpiredAt: m.ExpiredAt,
		CreatedAt: m.CreatedAt,
	}

	return resp, nil
}
