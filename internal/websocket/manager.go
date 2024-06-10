package websocket

import (
	"mt/config"
	"mt/internal/data"
	"mt/pkg/logger"
)

var (
	manager *Manager
)

func RegisterManagerInstance(managerInstance *Manager) { manager = managerInstance }

func ManagerInstance() *Manager { return manager }

type Manager struct {
	cServer  *config.Server
	resource *data.Data
	logger   *logger.Logger
}

func NewManager(cServer *config.Server, resource *data.Data, logger *logger.Logger) *Manager {
	return &Manager{
		cServer:  cServer,
		resource: resource,
		logger:   logger,
	}
}
