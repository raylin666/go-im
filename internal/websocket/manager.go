package websocket

import (
	"mt/internal/data"
	"mt/pkg/logger"
	"sync"
)

var (
	manager *Manager
)

func RegisterManager(m *Manager) { manager = m }

func ManagerInstance() *Manager { return manager }

type Manager struct {
	clientManager *ClientManager

	events       map[string]EventDisposeFunc
	eventRWMutex sync.RWMutex

	logger   *logger.Logger
	resource *data.Data
}

func NewManager(logger *logger.Logger, data *data.Data) (manager *Manager) {
	manager = &Manager{
		clientManager: NewClientManager(),
		events:        make(map[string]EventDisposeFunc),
		logger:        logger,
		resource:      data,
	}

	// 注册处理事件
	var events = NewEvents()
	for event, fn := range events.Registers {
		manager.WithEventHandler(event, fn)
	}

	return
}

func (m *Manager) Logger() *logger.Logger { return m.logger }

func (m *Manager) Resource() *data.Data { return m.resource }

func (m *Manager) ClientManager() *ClientManager { return m.clientManager }

// WithEventHandler 注册处理事件
func (m *Manager) WithEventHandler(event string, fn EventDisposeFunc) {
	m.eventRWMutex.Lock()
	defer m.eventRWMutex.Unlock()
	m.events[event] = fn
}

// WithEventHandler 获取处理事件
func (m *Manager) GetEventHandler(event string) (value EventDisposeFunc, ok bool) {
	m.eventRWMutex.RLock()
	defer m.eventRWMutex.RUnlock()
	value, ok = m.events[event]
	return
}
