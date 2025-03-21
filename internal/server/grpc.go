package server

import (
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"mt/api/v1"
	"mt/config"
	"mt/internal/app"
	"mt/internal/middleware/auth"
	logging "mt/internal/middleware/logger"
	"mt/internal/middleware/request"
	"mt/internal/middleware/validate"
	"mt/internal/service"

	accountPb "mt/api/v1/account"
	messagePb "mt/api/v1/message"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *config.Server,
	heartbeat *service.HeartbeatService,
	account *service.AccountService,
	message *service.MessageService,
	tools *app.Tools) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			metadata.Server(),
			request.Trace(),
			logging.Server(tools.Logger()),
			auth.NewJWTAuthServer(tools.JWT()),
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

	v1.RegisterHeartbeatServer(srv, heartbeat)
	accountPb.RegisterServiceServer(srv, account)
	messagePb.RegisterServiceServer(srv, message)

	return srv
}
