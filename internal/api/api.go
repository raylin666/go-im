package api

import (
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"mt/config"
	"mt/internal/app"
	"mt/internal/grpc"
	"mt/internal/websocket"
	"mt/pkg/repositories"
)

// ProviderSet is api.handler providers.
var ProviderSet = wire.NewSet(NewHandler)

type Handler struct {
	r               *mux.Router
	dbRepo          repositories.DbRepo
	redisRepo       repositories.RedisRepo
	grpcClient      *grpc.GrpcClient
	wsClientManager *websocket.ClientManager
	tools           *app.Tools
	config          *config.Bootstrap
	Prefix          string
}

func NewHandler(
	config *config.Bootstrap,
	tools *app.Tools,
	dataRepo repositories.DataRepo,
	grpcClient *grpc.GrpcClient,
	wsClientManager *websocket.ClientManager) *Handler {
	return &Handler{
		r:               mux.NewRouter(),
		dbRepo:          dataRepo.DbRepo(),
		redisRepo:       dataRepo.RedisRepo(),
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
