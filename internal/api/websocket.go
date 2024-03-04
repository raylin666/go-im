package api

import (
	"context"
	"go.uber.org/zap"
	"log"
	"mt/internal/constant/defined"
	"mt/pkg/websocket"
	"net/http"
)

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	var ctx = context.Background()
	c, err := websocket.NewUpgrader(w, r)
	if err != nil {
		var e = defined.ErrorWebsocketUpgraderError
		w.Write([]byte(e.GetReason()))
		w.WriteHeader(int(e.GetCode()))

		h.logger.UseApp(ctx).Error("WebSocket 连接失败", zap.Error(e))
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
