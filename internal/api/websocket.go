package api

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/google/uuid"
	"go.uber.org/zap"
	accountPb "mt/api/v1/account"
	"mt/internal/constant/defined"
	"mt/internal/lib"
	"mt/internal/websocket"
	"mt/pkg/logger"
	"mt/pkg/utils"
	"net/http"
)

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = metadata.NewServerContext(
			lib.NewContextHttpRequest(context.Background(), r),
			metadata.New(map[string][]string{logger.XMdKeyTraceId: {uuid.New().String()}}),
		)

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

	// 更新帐号登录信息
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
		h.tools.Logger().UseWebSocket(ctx).Error(fmt.Sprintf("WebSocket 建立连接失败: %s", conn.RemoteAddr().String()), zap.String("account_token", accountToken), zap.Any("account", account), zap.Error(err))
		return
	}

	h.tools.Logger().UseWebSocket(ctx).Info(fmt.Sprintf("WebSocket 建立连接完成: %s", conn.RemoteAddr().String()), zap.String("account_token", accountToken), zap.Any("account", account))

	// 创建客户端连接, 完成帐号连接信息存储
	client := h.wsClientManager.CreateClient(websocket.NewAccount(account.AccountId, account.Nickname, account.Avatar, account.IsAdmin), conn)

	// 监听客户端连接消息读写及事件处理
	h.wsClientManager.ClientRegister(client)
}
