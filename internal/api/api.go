package api

import (
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"mt/config"
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

	cApp       *config.App
	cWebsocket *config.Websocket
}

func NewHandler(cApp *config.App, cWebsocket *config.Websocket, logger *logger.Logger, dataRepo repositories.DataRepo) *Handler {
	return &Handler{
		r:         mux.NewRouter(),
		dbRepo:    dataRepo.DbRepo(),
		redisRepo: dataRepo.RedisRepo(),
		logger:    logger,
		Prefix:    "/app/",

		cApp:       cApp,
		cWebsocket: cWebsocket,
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
