package api

import (
	"context"
	"go.uber.org/zap"
	"mt/internal/constant/defined"
	"mt/internal/websocket"
	"net/http"
	"time"
)

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	var (
		timeNow = time.Now()

		ctx = context.TODO()

		query = r.URL.Query()
	)

	// TODO 登录身份验证
	accountToken := query.Get("account_token")
	if accountToken == "" {
		h.writeError(w, defined.ErrorNotVisitAuth)
		return
	}

	// TODO 解析TOKEN
	jwtClaims, err := h.tools.JWT().ParseToken(accountToken)
	if err != nil {
		h.writeError(w, defined.ErrorNotLoginError)
		return
	}

	// TODO HTTP 协议升级
	upgraderResponseHeader := new(websocket.UpgraderResponseHeader)
	upgraderResponseHeader.Name = h.config.App.Name
	upgraderResponseHeader.Version = h.config.App.Version
	conn, err := websocket.NewUpgrader(w, r, upgraderResponseHeader,
		websocket.WithUpgraderHandshakeTimeout(h.config.Websocket.HandshakeTimeout.AsDuration()),
		websocket.WithUpgraderReadBufferSize(int(h.config.Websocket.ReadBufferSize)),
		websocket.WithUpgraderWriteBufferSize(int(h.config.Websocket.WriteBufferSize)),
		websocket.WithUpgraderCheckOrigin(func(r *http.Request) bool {
			return true
		}),
		websocket.WithUpgraderError(func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			// TODO 升级失败处理
		}))
	if err != nil {
		var e = defined.ErrorWebsocketUpgraderError
		h.writeError(w, e)
		h.tools.Logger().UseWebSocket(ctx).Error("WebSocket 连接失败", zap.Error(e))
		return
	}
}
