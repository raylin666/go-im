package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/wire"
	"go.uber.org/zap"
	"mt/internal/app"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/query"
	"mt/internal/websocket"
	"mt/internal/websocket/types"
	"mt/pkg/repositories"
)

// ProviderSet is events providers.
var ProviderSet = wire.NewSet(NewEvents)

type events struct {
	relation map[string]websocket.EventDisposeFunc
	tools    *app.Tools
	repo     repositories.DataRepo
}

// NewEvents 消息事件
func NewEvents(tools *app.Tools, repo repositories.DataRepo) websocket.Events {
	return &events{tools: tools, repo: repo}
}

// GetAll 获取所有事件对应处理器
func (event *events) GetAll() map[string]websocket.EventDisposeFunc {
	event.relation = make(map[string]websocket.EventDisposeFunc)

	// 心跳检测处理器
	event.relation[websocket.EventPing] = event.Ping
	// 客户端和账号信息绑定
	event.relation[websocket.EventBind] = event.Bind
	// 获取账号信息
	event.relation[websocket.EventAccountInfo] = event.AccountInfo
	// 发送C2C消息
	event.relation[websocket.EventSendC2CMessage] = event.SendC2CMessage

	return event.relation
}

// ClientWhiteEventNames 客户端请求所支持的事件, 不在指定的事件客户端无法调用
func (event *events) ClientWhiteEventNames() []string {
	return []string{
		websocket.EventPing,
		websocket.EventBind,
		websocket.EventSendC2CMessage,
	}
}

// NewPushMessage 给客户端推送新消息事件
func (event *events) NewPushMessage(ctx context.Context, client *websocket.Client, eventName string, seq string, message []byte) bool {
	var loggerFields = event.loggerFields(eventName, seq, message)
	message, err := json.Marshal(&types.Request{Seq: seq, Event: eventName, Data: message})
	if err != nil {
		event.logger(ctx).Error(fmt.Sprintf("`%s` 消息推送给客户端失败", eventName), loggerFields...)
		return false
	}

	event.logger(ctx).Info(fmt.Sprintf("将 `%s` 消息推送给客户端", eventName), loggerFields...)
	client.EventMessageHandler(ctx, message, false)
	return true
}

// logger 获取日志对象
func (event *events) logger(ctx context.Context) *zap.Logger {
	return event.tools.Logger().UseWebSocket(ctx)
}

// dbQuery 获取数据库查询对象
func (event *events) dbQuery() *query.Query {
	return dbrepo.NewDefaultDbQuery(event.repo.DbRepo())
}

// loggerFields 获取事件日志字段信息
func (event *events) loggerFields(eventName, seq string, message []byte) []zap.Field {
	var logEvent = zap.String("event", eventName)
	var logSeq = zap.String("seq", seq)
	var logMessage = zap.String("message", string(message))

	return []zap.Field{logEvent, logSeq, logMessage}
}

// defaultEventResponse 默认事件返回值
func defaultEventResponse() (code uint32, msg string, data interface{}, send bool) {
	code, msg, data, send = 200, "OK", nil, true
	return
}
