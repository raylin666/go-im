package websocket

import (
	"fmt"
	"sync"
)

var (
	ClientManagerInstance *ClientManager
)

// ClientManager 连接管理
type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 APPID+UUID
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接处理
	// Login       chan *login        // 用户登录处理
	Unregister chan *Client // 断开连接处理
	Broadcast  chan []byte  // 广播消息-向全部成员发送数据
}

// NewClientManager 初始化连接管理
func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:  make(map[*Client]bool),
		Users:    make(map[string]*Client),
		Register: make(chan *Client, 1000),
		// Login:      make(chan *login, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}

	return
}

// GetUserKey 获取用户KEY
func GetUserKey(appId uint32, userId string) (key string) {
	key = fmt.Sprintf("%d_%s", appId, userId)
	return
}
