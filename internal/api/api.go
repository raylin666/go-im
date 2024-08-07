package api

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"mt/config"
	"mt/internal/app"
	"mt/pkg/repositories"
)

// ProviderSet is api.handler providers.
var ProviderSet = wire.NewSet(NewHandler)

type Handler struct {
	r         *mux.Router
	dbRepo    repositories.DbRepo
	redisRepo repositories.RedisRepo
	tools     *app.Tools
	Prefix    string

	cApp       *config.App
	cWebsocket *config.Websocket
}

func NewHandler(cApp *config.App, cWebsocket *config.Websocket, tools *app.Tools, dataRepo repositories.DataRepo) *Handler {
	return &Handler{
		r:         mux.NewRouter(),
		dbRepo:    dataRepo.DbRepo(),
		redisRepo: dataRepo.RedisRepo(),
		tools:     tools,
		Prefix:    "/app/",

		cApp:       cApp,
		cWebsocket: cWebsocket,
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

// writeError 写入错误信息
func (h *Handler) writeError(w http.ResponseWriter, err *errors.Error) {
	_, _ = w.Write([]byte(err.Reason))
	w.WriteHeader(int(err.Code))
	return
}
