package service

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"mt/api/v1"
	"mt/internal/biz"
	"mt/internal/constant/types"
)

type ManagerService struct {
	v1.UnimplementedManagerServer

	uc *biz.ManagerUsecase
}

func NewManagerService(uc *biz.ManagerUsecase) *ManagerService {
	return &ManagerService{uc: uc}
}

// Create 创建应用
func (s *ManagerService) Create(ctx context.Context, req *v1.ManagerCreateRequest) (*v1.ManagerCreateResponse, error) {
	createRequest := &types.CreateRequest{
		Ident:     req.GetIdent(),
		Name:      req.GetName(),
		Status:    uint32(req.GetStatus()),
		ExpiredAt: req.GetExpiredAt().AsTime(),
	}

	createResponse, err := s.uc.Create(ctx, createRequest)
	if err != nil {
		return nil, err
	}

	resp := &v1.ManagerCreateResponse{
		Id:        uint64(createResponse.Id),
		Ident:     createResponse.Ident,
		Name:      createResponse.Name,
		Key:       uint64(createResponse.Key),
		Secret:    createResponse.Secret,
		Status:    v1.ManagerStatus(createResponse.Status),
		ExpiredAt: timestamppb.New(createResponse.ExpiredAt),
		CreatedAt: timestamppb.New(createResponse.CreatedAt),
	}

	return resp, nil
}
