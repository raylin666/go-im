package websocket

import (
	"mt/config"
	"mt/internal/app"
	"mt/internal/data"
	"sync"
)

var (
	manager *Manager
)

func RegisterManagerInstance(managerInstance *Manager) { manager = managerInstance }

func ManagerInstance() *Manager { return manager }

// Manager WebSocket 资源分发管理器
type Manager struct {
	clientManager *ClientManager
	events        map[string]EventDisposeFunc
	eventRWMutex  sync.RWMutex
	cServer       *config.Server
	resource      *data.Data
	tools         *app.Tools
}

func NewManager(cServer *config.Server, resource *data.Data, tools *app.Tools, events Events) (manager *Manager) {
	manager = &Manager{
		clientManager: NewClientManager(),
		events:        make(map[string]EventDisposeFunc),
		cServer:       cServer,
		resource:      resource,
		tools:         tools,
	}

	// 注册消息事件处理器
	for event, fn := range events.GetAll() {
		manager.WithEventHandler(event, fn)
	}

	// 注册管道事件处理
	go manager.clientManager.ChanEventHandler()

	return
}

func (manager *Manager) ClientManager() *ClientManager { return manager.clientManager }

// WithEventHandler 注册消息事件处理器
func (manager *Manager) WithEventHandler(event string, fn EventDisposeFunc) {
	manager.eventRWMutex.Lock()
	defer manager.eventRWMutex.Unlock()
	manager.events[event] = fn
}

// GetEventHandler 获取消息事件处理器
func (manager *Manager) GetEventHandler(event string) (value EventDisposeFunc, ok bool) {
	manager.eventRWMutex.RLock()
	defer manager.eventRWMutex.RUnlock()
	value, ok = manager.events[event]
	return
}
