package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/constant/defined"
	"mt/internal/lib"
	"mt/internal/repositories/dbrepo"
	"mt/internal/websocket"
	"net/http"
)

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = lib.NewContextHttpRequest(context.Background(), r)

		query = r.URL.Query()
	)

	// TODO 登录身份验证
	accountToken := query.Get("account_token")
	if accountToken == "" {
		var e = defined.ErrorNotVisitAuth
		_, _ = w.Write([]byte(e.Reason))
		w.WriteHeader(int(e.Code))

		return
	}

	// TODO 解析TOKEN
	jwtClaims, err := h.tools.JWT().ParseToken(accountToken)
	if err != nil {
		var e = defined.ErrorAuthenticationError
		_, _ = w.Write([]byte(e.Reason))
		w.WriteHeader(int(e.Code))

		return
	}

	// TODO 账号校验
	q := dbrepo.NewDefaultDbQuery(h.dbRepo).Account
	account, err := q.WithContext(ctx).FirstByAccountId(jwtClaims.ID)
	if err != nil {
		var e = defined.ErrorAccountLoginError
		_, _ = w.Write([]byte(e.Reason))
		w.WriteHeader(int(e.Code))

		return
	}

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
			// TODO 升级失败处理
		}))
	if err != nil {
		var e = defined.ErrorWebsocketUpgraderError
		_, _ = w.Write([]byte(e.Reason))
		w.WriteHeader(int(e.Code))

		h.tools.Logger().UseWebSocket(ctx).Error("WebSocket 连接失败", zap.Error(e))
		return
	}

	h.tools.Logger().UseWebSocket(ctx).Info(fmt.Sprintf("WebSocket 建立连接: %s", conn.RemoteAddr().String()))

	client := websocket.NewClient(account.AccountId, conn)

	go client.Read(ctx)
	go client.Write(ctx)

	websocket.ManagerInstance().ClientManager().Register <- client
}
