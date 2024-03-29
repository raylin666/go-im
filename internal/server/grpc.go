package server

import (
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"mt/api/v1"
	accountPb "mt/api/v1/account"
	managerPb "mt/api/v1/manager"
	"mt/config"
	"mt/internal/middleware/auth"
	logging "mt/internal/middleware/logger"
	"mt/internal/middleware/validate"
	"mt/internal/service"
	"mt/pkg/logger"
	"mt/pkg/repositories"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *config.Server,
	apiHeartbeat *service.HeartbeatService,
	apiManager *service.ManagerService,
	apiAccount *service.AccountService,
	repo repositories.DataRepo,
	logger *logger.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			metadata.Server(),
			logging.Server(logger),
			auth.NewJWTAuthServer(),
			auth.NewAccountAuthServer(repo),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)

	v1.RegisterHeartbeatServer(srv, apiHeartbeat)
	managerPb.RegisterServiceServer(srv, apiManager)
	accountPb.RegisterServiceServer(srv, apiAccount)

	return srv
}
