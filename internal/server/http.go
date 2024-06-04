package server

import (
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"mt/api/v1"
	"mt/config"
	"mt/internal/api"
	"mt/internal/middleware/auth"
	"mt/internal/middleware/encode"
	logging "mt/internal/middleware/logger"
	"mt/internal/service"
	"mt/pkg/logger"
	netHttp "net/http"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	accountPb "mt/api/v1/account"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *config.Server,
	heartbeat *service.HeartbeatService,
	account *service.AccountService,
	apiHandler *api.Handler,
	logger *logger.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			metadata.Server(),
			logging.Server(logger),
			auth.NewJWTAuthServer(),
			// response.ErrorHandler(),
		),
		http.ResponseEncoder(encode.ResponseEncoder),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	// HTTP API 路由处理器
	srv.HandlePrefix(apiHandler.Prefix, netHttp.Handler(apiHandler.Router()))

	v1.RegisterHeartbeatHTTPServer(srv, heartbeat)
	accountPb.RegisterServiceHTTPServer(srv, account)

	return srv
}
