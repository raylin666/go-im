package websocket

import (
	"mt/internal/data"
	"mt/pkg/logger"
)

var (
	manager *Manager
)

func RegisterManager(m *Manager) { manager = m }

func ManagerInstance() *Manager { return manager }

type Manager struct {
	clientManager *ClientManager

	logger   *logger.Logger
	resource *data.Data
}

func NewManager(logger *logger.Logger, data *data.Data) *Manager {
	return &Manager{
		clientManager: NewClientManager(),
		logger:        logger,
		resource:      data,
	}
}

func (m *Manager) Logger() *logger.Logger {
	return m.logger
}

func (m *Manager) Resource() *data.Data {
	return m.resource
}

func (m *Manager) ClientManager() *ClientManager {
	return m.clientManager
}
