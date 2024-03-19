package service

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "mt/api/v1/manager"
	"mt/internal/biz"
	"mt/internal/constant/types"
)

type ManagerService struct {
	pb.UnimplementedManagerServer

	uc *biz.ManagerUsecase
}

func NewManagerService(uc *biz.ManagerUsecase) *ManagerService {
	return &ManagerService{uc: uc}
}

// Create 创建应用
func (s *ManagerService) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	createRequest := &types.ManagerCreateRequest{
		Ident:     req.GetIdent(),
		Name:      req.GetName(),
		Status:    int8(req.GetStatus()),
		ExpiredAt: req.GetExpiredAt().AsTime(),
	}

	createResponse, err := s.uc.Create(ctx, createRequest)
	if err != nil {
		return nil, err
	}

	resp := &pb.CreateResponse{
		Id:        uint64(createResponse.Id),
		Ident:     createResponse.Ident,
		Name:      createResponse.Name,
		Key:       createResponse.Key,
		Secret:    createResponse.Secret,
		Status:    pb.Status(createResponse.Status),
		ExpiredAt: timestamppb.New(createResponse.ExpiredAt),
		CreatedAt: timestamppb.New(createResponse.CreatedAt),
	}

	return resp, nil
}
