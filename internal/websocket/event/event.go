package event

import (
	"context"
	"github.com/google/wire"
	"go.uber.org/zap"
	"mt/internal/app"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/query"
	"mt/internal/websocket"
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

	return event.relation
}

// logger 获取日志对象
func (event *events) logger(ctx context.Context) *zap.Logger {
	return event.tools.Logger().UseWebSocket(ctx)
}

// dbQuery 获取数据库查询对象
func (event *events) dbQuery() *query.Query {
	return dbrepo.NewDefaultDbQuery(event.repo.DbRepo())
}

// defaultEventResponse 默认事件返回值
func defaultEventResponse() (code uint32, msg string, data interface{}, send bool) {
	code, msg, data, send = 200, "OK", nil, true
	return
}
