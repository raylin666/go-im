package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/constant/defined"
	"mt/internal/lib"
	"mt/internal/websocket"
	"net/http"
)

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = lib.NewContextHttpRequest(context.Background(), r)

		//query = r.URL.Query()
	)

	// TODO HTTP 协议升级
	upgraderResponseHeader := new(websocket.UpgraderResponseHeader)
	upgraderResponseHeader.Name = h.cApp.Name
	upgraderResponseHeader.Version = h.cApp.Version
	conn, err := websocket.NewUpgrader(w, r, upgraderResponseHeader,
		websocket.WithUpgraderHandshakeTimeout(h.cWebsocket.HandshakeTimeout.AsDuration()),
		websocket.WithUpgraderReadBufferSize(int(h.cWebsocket.ReadBufferSize)),
		websocket.WithUpgraderWriteBufferSize(int(h.cWebsocket.WriteBufferSize)),
		websocket.WithUpgraderCheckOrigin(func(r *http.Request) bool {
			return true
		}),
		websocket.WithUpgraderError(func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			fmt.Println(status, reason)
		}))
	if err != nil {
		var e = defined.ErrorWebsocketUpgraderError
		_, _ = w.Write([]byte(e.Reason))
		w.WriteHeader(int(e.Code))

		h.logger.UseWebSocket(ctx).Error("WebSocket 连接失败", zap.Error(e))
		return
	}

	h.logger.UseWebSocket(ctx).Info(fmt.Sprintf("WebSocket 建立连接: %s", conn.RemoteAddr().String()))
}
