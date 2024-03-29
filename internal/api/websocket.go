package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/constant/defined"
	"mt/internal/lib"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/model"
	"mt/internal/websocket"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = lib.NewContextHttpRequest(context.Background(), r)

		query = r.URL.Query()
	)

	// TODO 应用身份验证
	appKey, err := strconv.Atoi(query.Get("key"))
	appSecret := query.Get("secret")
	if err != nil || appKey <= 0 || appSecret == "" {
		var e = defined.ErrorNotVisitAuth
		_, _ = w.Write([]byte(e.GetReason()))
		w.WriteHeader(int(e.GetCode()))

		return
	}

	q := dbrepo.NewDefaultDbQuery(h.dbRepo).App
	m, _ := q.WithContext(ctx).FirstByKeyAndSecret(uint64(appKey), appSecret)
	modelErr := model.AppAvailableByKeyAndSecret(&m)
	if modelErr != nil {
		_, _ = w.Write([]byte(modelErr.GetReason()))
		w.WriteHeader(int(modelErr.GetCode()))

		return
	}

	// TODO HTTP 协议升级
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

		h.logger.UseWebSocket(ctx).Error("WebSocket 连接失败", zap.Error(e))
		return
	}

	h.logger.UseWebSocket(ctx).Info(fmt.Sprintf("WebSocket 建立连接: %s", conn.RemoteAddr().String()))

	client := websocket.NewClient(uint64(appKey), conn)

	go client.Read(ctx)
	go client.Write(ctx)

	// 用户连接处理
	websocket.ManagerInstance().ClientManager().Register <- client
}
