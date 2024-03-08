package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/constant/defined"
	websocket2 "mt/internal/websocket"
	"net/http"
)

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()

	conn, err := websocket2.NewUpgrader(w, r)
	if err != nil {
		var e = defined.ErrorWebsocketUpgraderError
		_, _ = w.Write([]byte(e.GetReason()))
		w.WriteHeader(int(e.GetCode()))

		h.logger.UseApp(ctx).Error("WebSocket 连接失败", zap.Error(e))
		return
	}

	h.logger.UseApp(ctx).Info(fmt.Sprintf("WebSocket 建立连接: %s", conn.RemoteAddr().String()))

	client := websocket2.NewClient(ctx, h.logger, conn)

	go client.Read()
	go client.Write()

	// 用户连接事件
	websocket2.ClientManagerInstance.Register <- client
}
