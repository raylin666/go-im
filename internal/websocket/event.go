package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/constant/defined"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/model"
	typesEvent "mt/internal/websocket/types/event"
	"mt/pkg/utils"
	"time"
)

type EventDisposeFunc func(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

type Events struct {
	Registers map[string]EventDisposeFunc
}

func NewEvents() (events *Events) {
	events = &Events{}
	events.Registers = make(map[string]EventDisposeFunc)

	// 注册处理事件
	events.Registers["ping"] = events.Ping
	events.Registers["login"] = events.Login

	return
}

// Ping 心跳检测
func (event *Events) Ping(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = codeStatusOk
	msg = codeMessageOk
	data = "pong"

	client.Heartbeat(uint64(time.Now().Unix()))

	return
}

// Login 账号登录 (必须登录完成后才能进行用户事件)
func (event *Events) Login(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = codeStatusOk
	msg = codeMessageOk
	data = nil

	request := &typesEvent.LoginRequest{}
	err := json.Unmarshal(message, request)
	if err != nil || request.UserId == "" {
		Logger(ctx).Error("账号登录事件-解析消息数据包错误 json.Marshal", zap.Error(err))
		code, msg = utils.ErrorMessage(defined.ErrorRequestParamsError)
		return
	}

	q := dbrepo.NewDefaultDbQuery(DbRepo()).Account.Table(model.AccountTableName(client.AppKey))
	m, err := q.WithContext(ctx).FirstByUserId(request.UserId)
	fmt.Println(m)

	return
}
