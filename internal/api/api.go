package api

import (
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"mt/config"
	"mt/internal/app"
	"mt/internal/grpc"
	"mt/internal/repositories"
	"mt/internal/websocket"
)

// ProviderSet is api.handler providers.
var ProviderSet = wire.NewSet(NewHandler)

type Handler struct {
	r               *mux.Router
	dataRepo        repositories.DataRepo
	grpcClient      grpc.GrpcClient
	wsClientManager websocket.WebsocketClientManager
	tools           *app.Tools
	config          *config.Bootstrap
	Prefix          string
}

func NewHandler(
	config *config.Bootstrap,
	tools *app.Tools,
	dataRepo repositories.DataRepo,
	grpcClient grpc.GrpcClient,
	wsClientManager websocket.WebsocketClientManager) *Handler {
	return &Handler{
		r:               mux.NewRouter(),
		dataRepo:        dataRepo,
		grpcClient:      grpcClient,
		wsClientManager: wsClientManager,
		tools:           tools,
		config:          config,
		Prefix:          "/app/",
	}
}

func (h *Handler) Router() *mux.Router {
	// WebSocket
	var ws = h.r.PathPrefix(h.Prefix + "ws").Subrouter()
	{
		ws.HandleFunc("", h.WebSocket)
	}

	return h.r
}
