package event

import (
	"context"
	"github.com/google/wire"
	"mt/internal/websocket"
	"time"
)

// ProviderSet is events providers.
var ProviderSet = wire.NewSet(NewEventRepository)

type EventRepository struct {
	events Events
}

type Events interface {
	Ping(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{})
}

func NewEventRepository(events Events) *EventRepository {
	return &EventRepository{events: events}
}

// Ping 心跳检测[消息事件处理]
func (event EventRepository) Ping(ctx context.Context, client *websocket.Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code, msg, _ = defaultEventResponse()
	data = "pong"

	client.Heartbeat(time.Now())

	return
}

// defaultEventResponse 默认事件返回值
func defaultEventResponse() (code uint32, msg string, data interface{}) {
	code, msg, data = 200, "OK", nil
	return
}
