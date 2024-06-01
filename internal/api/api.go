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
		Prefix:    "/app/",
	}
}

func (h *Handler) Router() *mux.Router {
	// WebSocket
	var ws = h.r.PathPrefix(h.Prefix + "ws").Subrouter()
	h.routerWS(ws)

	return h.r
}

func (h *Handler) routerWS(r *mux.Router) {
	r.HandleFunc("", h.WebSocket)
}
