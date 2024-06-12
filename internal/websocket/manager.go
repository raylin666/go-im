package websocket

import (
	"mt/config"
	"mt/internal/data"
	"mt/internal/event"
	"mt/pkg/logger"
	"sync"
)

var (
	manager *Manager
)

func RegisterManagerInstance(managerInstance *Manager) {
	manager = managerInstance
}

func ManagerInstance() *Manager { return manager }

// Manager WebSocket 资源分发管理器
type Manager struct {
	clientManager *ClientManager

	events       map[string]EventDisposeFunc
	eventRWMutex sync.RWMutex

	cServer  *config.Server
	resource *data.Data
	logger   *logger.Logger
}

func NewManager(cServer *config.Server, resource *data.Data, logger *logger.Logger, events event.Events) (manager *Manager) {
	manager = &Manager{
		clientManager: NewClientManager(),
		events:        make(map[string]EventDisposeFunc),
		cServer:       cServer,
		resource:      resource,
		logger:        logger,
	}

	// 注册消息事件处理
	for event, fn := range Events() {
		manager.WithEventHandler(event, fn)
	}

	// 注册管道事件处理
	go manager.clientManager.ChanEventHandler()

	return
}

func (manager *Manager) ClientManager() *ClientManager { return manager.clientManager }

// WithEventHandler 注册消息事件处理
func (m *Manager) WithEventHandler(event string, fn EventDisposeFunc) {
	m.eventRWMutex.Lock()
	defer m.eventRWMutex.Unlock()
	m.events[event] = fn
}

// GetEventHandler 获取消息事件处理
func (m *Manager) GetEventHandler(event string) (value EventDisposeFunc, ok bool) {
	m.eventRWMutex.RLock()
	defer m.eventRWMutex.RUnlock()
	value, ok = m.events[event]
	return
}
