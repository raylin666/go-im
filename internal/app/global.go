package app

import (
	"github.com/raylin666/go-utils/auth"
	"github.com/raylin666/go-utils/server/system"
	"mt/pkg/websocket"
)

var (
	Datetime *system.Datetime
	JWT      auth.JWT

	// ClientManager 客户端连接管理
	ClientManager *websocket.ClientManager
)
