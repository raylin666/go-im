package server

import (
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"mt/api/v1"
	"mt/config"
	"mt/internal/middleware/auth"
	logging "mt/internal/middleware/logger"
	"mt/internal/middleware/validate"
	"mt/internal/service"
	"mt/pkg/logger"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	accountPb "mt/api/v1/account"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(cServer *config.Server,
	heartbeat *service.HeartbeatService,
	account *service.AccountService,
	logger *logger.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			metadata.Server(),
			logging.Server(logger),
			auth.NewJWTAuthServer(),
		),
	}
	if cServer.Grpc.Network != "" {
		opts = append(opts, grpc.Network(cServer.Grpc.Network))
	}
	if cServer.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(cServer.Grpc.Addr))
	}
	if cServer.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(cServer.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)

	v1.RegisterHeartbeatServer(srv, heartbeat)
	accountPb.RegisterServiceServer(srv, account)

	return srv
}
