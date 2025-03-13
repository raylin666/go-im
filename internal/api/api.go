package api

import (
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"mt/config"
	"mt/internal/app"
	"mt/internal/data"
	"mt/internal/grpc"
	"mt/internal/websocket"
)

// ProviderSet is api.handler providers.
var ProviderSet = wire.NewSet(NewHandler)

type Handler struct {
	r               *mux.Router
	grpcClient      grpc.GrpcClient
	wsClientManager websocket.WSClientManager
	tools           *app.Tools
	config          *config.Bootstrap
	Prefix          string

	// 数据逻辑仓库
	dataLogicRepo struct {
		Account data.AccountRepo
	}
}

func NewHandler(
	config *config.Bootstrap,
	tools *app.Tools,
	grpcClient grpc.GrpcClient,
	wsClientManager websocket.WSClientManager,
	accountRepo data.AccountRepo) *Handler {
	var handler = &Handler{
		r:               mux.NewRouter(),
		grpcClient:      grpcClient,
		wsClientManager: wsClientManager,
		tools:           tools,
		config:          config,
		Prefix:          "/app/",
	}

	handler.dataLogicRepo.Account = accountRepo
	return handler
}

func (h *Handler) Router() *mux.Router {
	// WebSocket
	var ws = h.r.PathPrefix(h.Prefix + "ws").Subrouter()
	{
		ws.HandleFunc("", h.WebSocket)
	}

	return h.r
}
