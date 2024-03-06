package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
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
		_, _ = w.Write([]byte(e.GetReason()))
		w.WriteHeader(int(e.GetCode()))

		h.logger.UseApp(ctx).Error("WebSocket 连接失败", zap.Error(e))
		return
	}

	currentTime := uint64(time.Now().Unix())
	connRemoteAddr := conn.RemoteAddr().String()

	h.logger.UseApp(ctx).Info(fmt.Sprintf("WebSocket 建立连接: %s", connRemoteAddr))

	client := websocket.NewClient(ctx, h.logger, connRemoteAddr, conn, currentTime)

	go client.Read()
	go client.Write()

	// 用户连接事件
	websocket.ClientManagerInstance.Register <- client
}
