package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/app"
	"mt/internal/constant/defined"
	"mt/pkg/websocket"
	"net/http"
	"time"
)

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()
	conn, err := websocket.NewUpgrader(w, r)
	if err != nil {
		var e = defined.ErrorWebsocketUpgraderError
		w.Write([]byte(e.GetReason()))
		w.WriteHeader(int(e.GetCode()))

		h.logger.UseApp(ctx).Error("WebSocket 连接失败", zap.Error(e))
		return
	}

	h.logger.UseApp(ctx).Info(fmt.Sprintf("WebSocket 建立连接: %s", conn.RemoteAddr().String()))

	currentTime := uint64(time.Now().Unix())
	client := websocket.NewClient(conn.RemoteAddr().String(), conn, currentTime)

	go client.Read()
	go client.Write()

	// 用户连接事件
	app.ClientManager.Register <- client
}
