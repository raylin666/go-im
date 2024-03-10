package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/constant/defined"
	"mt/internal/websocket"
	"net/http"
	"time"
)

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()

	upgraderResponseHeader := new(websocket.UpgraderResponseHeader)
	upgraderResponseHeader.Name = "goim"
	upgraderResponseHeader.Version = "1.0"
	conn, err := websocket.NewUpgrader(w, r, upgraderResponseHeader,
		websocket.WithUpgraderHandshakeTimeout(5*time.Second),
		websocket.WithUpgraderReadBufferSize(2048),
		websocket.WithUpgraderWriteBufferSize(2048),
		websocket.WithUpgraderCheckOrigin(func(r *http.Request) bool {
			return true
		}),
		websocket.WithUpgraderError(func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			fmt.Println(status, reason)
		}))
	if err != nil {
		var e = defined.ErrorWebsocketUpgraderError
		_, _ = w.Write([]byte(e.GetReason()))
		w.WriteHeader(int(e.GetCode()))

		h.logger.UseApp(ctx).Error("WebSocket 连接失败", zap.Error(e))
		return
	}

	h.logger.UseApp(ctx).Info(fmt.Sprintf("WebSocket 建立连接: %s", conn.RemoteAddr().String()))

	client := websocket.NewClient(conn)

	go client.Read(ctx)
	go client.Write(ctx)

	// 用户连接处理
	websocket.ManagerInstance().ClientManager().Register <- client
}
