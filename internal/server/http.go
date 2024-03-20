package server

import (
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"mt/api/v1"
	accountPb "mt/api/v1/account"
	managerPb "mt/api/v1/manager"
	"mt/config"
	"mt/internal/api"
	"mt/internal/middleware/auth"
	"mt/internal/middleware/encode"
	logging "mt/internal/middleware/logger"
	"mt/internal/service"
	"mt/internal/websocket"
	"mt/pkg/logger"
	netHttp "net/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(
	c *config.Server,
	apiHeartbeat *service.HeartbeatService,
	apiManager *service.ManagerService,
	apiAccount *service.AccountService,
	apiHandler *api.Handler,
	wsManager *websocket.Manager,
	logger *logger.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			metadata.Server(),
			logging.Server(logger),
			auth.NewJWTAuthServer(),
			auth.NewAccountAuthServer(),
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

	// 注册 WebSocket 管理器
	websocket.RegisterManager(wsManager)

	// HTTP API 路由处理器
	srv.HandlePrefix(apiHandler.Prefix, netHttp.Handler(apiHandler.Router()))

	// HTTP 服务路由处理器
	v1.RegisterHeartbeatHTTPServer(srv, apiHeartbeat)
	managerPb.RegisterManagerHTTPServer(srv, apiManager)
	accountPb.RegisterAccountHTTPServer(srv, apiAccount)

	return srv
}
