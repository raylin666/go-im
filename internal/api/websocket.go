package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	accountPb "mt/api/v1/account"
	"mt/internal/constant/defined"
	"mt/internal/lib"
	"mt/internal/websocket"
	"mt/pkg/utils"
	"net/http"
)

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = lib.NewContextHttpRequest(context.Background(), r)

		query = r.URL.Query()

		clientIp = utils.ClientIP(r)
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

	account, err := h.grpcClient.Account.UpdateLogin(ctx, &accountPb.UpdateLoginRequest{AccountId: jwtClaims.ID, ClientIp: clientIp})
	if err != nil {
		h.writeError(w, defined.ErrorAccountLoginError)
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

	h.tools.Logger().UseWebSocket(ctx).Info(fmt.Sprintf("WebSocket 建立连接: %s", conn.RemoteAddr().String()), zap.String("account_token", accountToken), zap.Any("account", account))

	// client := websocket.NewClient(ctx, websocket.NewAccount(account.AccountId, account.Nickname, account.Avatar, account.IsAdmin == 1), conn)

	// go client.Read(ctx)
	// go client.Write(ctx)
}
