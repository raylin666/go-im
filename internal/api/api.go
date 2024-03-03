package api

import (
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"mt/pkg/logger"
	"mt/pkg/repositories"
)

// ProviderSet is api.handler providers.
var ProviderSet = wire.NewSet(NewHandler)

type Handler struct {
	r         *mux.Router
	dbRepo    repositories.DbRepo
	redisRepo repositories.RedisRepo
	logger    *logger.Logger
	Prefix    string
}

func NewHandler(logger *logger.Logger, dataRepo repositories.DataRepo) *Handler {
	return &Handler{
		r:         mux.NewRouter(),
		dbRepo:    dataRepo.DbRepo(),
		redisRepo: dataRepo.RedisRepo(),
		logger:    logger,
		Prefix:    "/api/",
	}
}

func (h *Handler) Router() *mux.Router {
	// HTTP
	h.r.HandleFunc(h.Prefix+"account", h.Account).Methods("GET")

	// WebSocket
	h.r.HandleFunc("/ws", h.WebSocket)

	return h.r
}
