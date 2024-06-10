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
	"mt/internal/websocket"
	"mt/pkg/logger"
	netHttp "net/http"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	accountPb "mt/api/v1/account"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(cServer *config.Server,
	heartbeat *service.HeartbeatService,
	account *service.AccountService,
	apiHandler *api.Handler,
	websocketManager *websocket.Manager,
	logger *logger.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			metadata.Server(),
			logging.Server(logger),
			auth.NewJWTAuthServer(),
		),
		// 请求响应序列化
		http.ResponseEncoder(encode.ResponseEncoder),
		// 错误响应序列化
		http.ErrorEncoder(encode.ErrorEncoder),
	}
	if cServer.Http.Network != "" {
		opts = append(opts, http.Network(cServer.Http.Network))
	}
	if cServer.Http.Addr != "" {
		opts = append(opts, http.Address(cServer.Http.Addr))
	}
	if cServer.Http.Timeout != nil {
		opts = append(opts, http.Timeout(cServer.Http.Timeout.AsDuration()))
	}

	// 注册 Websocket 管理器
	websocket.RegisterManagerInstance(websocketManager)

	srv := http.NewServer(opts...)

	// HTTP API 路由处理器
	srv.HandlePrefix(apiHandler.Prefix, netHttp.Handler(apiHandler.Router()))

	v1.RegisterHeartbeatHTTPServer(srv, heartbeat)
	accountPb.RegisterServiceHTTPServer(srv, account)

	return srv
}
