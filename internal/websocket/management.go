package websocket

import (
	"mt/internal/app"
	"mt/pkg/repositories"
)

type Management struct {
	tools    *app.Tools
	dataRepo repositories.DataRepo

	ClientManager *ClientManager
}

func NewManagement(dataRepo repositories.DataRepo, tools *app.Tools) *Management {
	var management = &Management{
		tools:    tools,
		dataRepo: dataRepo,

		ClientManager: NewClientManager(),
	}

	return management
}
