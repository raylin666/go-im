package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gen/field"
	"mt/internal/constant/defined"
	"mt/internal/lib"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/model"
	"mt/internal/websocket"
	"mt/pkg/utils"
	"net/http"
	"time"
)

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	var (
		timeNow = time.Now()

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

	// TODO 账号校验
	accountQuery := dbrepo.NewDefaultDbQuery(h.dbRepo).Account
	account, err := accountQuery.WithContext(ctx).FirstByAccountId(jwtClaims.ID)
	if err != nil {
		h.writeError(w, defined.ErrorAccountLoginError)
		return
	}

	// TODO 处理账号登录, 更新账号信息
	var assignExpr = []field.AssignExpr{
		accountQuery.Status.Value(model.AccountStatusOnline),
		accountQuery.LastLoginTime.Value(timeNow),
		accountQuery.LastLoginIp.Value(clientIp),
		accountQuery.UpdatedAt.Value(timeNow),
	}

	accountOnlineQuery := dbrepo.NewDefaultDbQuery(h.dbRepo).AccountOnline
	if accountOnlineExistsResult, err := accountOnlineQuery.WithContext(ctx).ExistsByAccountId(account.AccountId); err == nil {
		if existsResult, existsResultOk := accountOnlineExistsResult["ok"]; existsResultOk {
			existsValue, existsValueOk := existsResult.(int64)
			if existsValueOk && existsValue == 0 {
				assignExpr = append(assignExpr, accountQuery.FirstLoginTime.Value(timeNow))
			}
		}
	}

	_, err = accountQuery.WithContext(ctx).Where(accountQuery.AccountId.Eq(account.AccountId)).UpdateSimple(assignExpr...)
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

	h.tools.Logger().UseWebSocket(ctx).Info(fmt.Sprintf("WebSocket 建立连接: %s", conn.RemoteAddr().String()), zap.String("account_token", accountToken), zap.Any("account", account))

	client := websocket.NewClient(ctx, websocket.NewAccount(account.AccountId, account.Nickname, account.Avatar, account.IsAdmin == 1), conn)

	go client.Read(ctx)
	go client.Write(ctx)
}
