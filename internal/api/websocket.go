package api

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/google/uuid"
	gorillaWebsocket "github.com/gorilla/websocket"
	"go.uber.org/zap"
	"mt/errors"
	"mt/internal/constant/types"
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
	)

	// TODO 登录身份验证
	accountToken := query.Get("account_token")
	if accountToken == "" {
		e := errors.New().NotVisitAuth()
		http.Error(w, e.GetMessage(), int(e.GetCode()))
		return
	}

	// TODO 解析TOKEN
	jwtClaims, err := h.tools.JWT().ParseToken(accountToken)
	if err != nil {
		e := errors.New().NotLogin()
		http.Error(w, e.GetMessage(), int(e.GetCode()))
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
		var e = errors.New().ServerUpgrader()
		http.Error(w, e.GetMessage(), int(e.GetCode()))
		h.tools.Logger().UseWebSocket(ctx).Error(fmt.Sprintf("WebSocket 建立连接失败: %s", conn.RemoteAddr().String()), zap.String("account_token", accountToken), zap.Any("account_id", jwtClaims.ID), zap.Error(err))
		return
	}

	// TODO 登录帐号
	account, accountOnline, err := h.dataLogicRepo.Account.Login(ctx, jwtClaims.ID, &types.AccountLoginRequest{
		ClientIp:   utils.ClientIP(r),
		ClientAddr: conn.LocalAddr().String(),
		ServerAddr: conn.RemoteAddr().String(),
		DeviceId:   r.Header.Get("device_id"),
		Os:         r.Header.Get("os"),
		System:     r.Header.Get("system"),
	})

	if err != nil {
		conn.WriteMessage(gorillaWebsocket.TextMessage, []byte(fmt.Sprintf("%s: %s", errors.New().AccountLogin().GetMessage(), err.Error())))
		conn.Close()
		return
	}

	h.tools.Logger().UseWebSocket(ctx).Info(fmt.Sprintf("WebSocket 建立连接完成: %s", conn.RemoteAddr().String()), zap.String("account_token", accountToken), zap.Any("account", account), zap.Any("account_online", accountOnline))

	// TODO 创建客户端连接, 完成帐号连接信息存储
	client := h.wsClientManager.CreateClient(ctx, websocket.NewAccount(account.AccountId, account.Nickname, account.Avatar, accountOnline.ID, account.IsAdmin == 1), conn)

	// TODO 监听客户端连接消息读写及事件处理
	h.wsClientManager.ClientRegister(client)
}
