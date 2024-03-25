package websocket

import (
	"go.uber.org/zap"
	"sync"
	"time"
)

// ClientManager 连接管理
type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 APPID_UUID
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 建立连接处理
	UnRegister  chan *Client       // 断开连接处理
	Broadcast   chan []byte        // 广播消息-向全部成员发送数据
}

// NewClientManager 初始化连接管理
func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 1000),
		UnRegister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}

	return
}

// EventRegister 建立连接事件
func (manager *ClientManager) EventRegister(client *Client) {
	manager.CreateClient(client)

	Logger(ctx).Info("客户端管理器 - 建立连接事件",
		zap.String("address", client.Addr),
		zap.Uint64("app_key", client.AppKey),
		zap.Time("heartbeat_time", time.Unix(int64(client.HeartbeatTime), 0)))
}

// EventUnRegister 断开连接事件
func (manager *ClientManager) EventUnRegister(client *Client) {
	Logger(ctx).Info("客户端管理器 - 断开连接事件",
		zap.String("address", client.Addr),
		zap.Uint64("app_key", client.AppKey),
		zap.Time("first_time", time.Unix(int64(client.FirstTime), 0)),
		zap.Time("heartbeat_time", time.Unix(int64(client.HeartbeatTime), 0)),
		zap.String("user_id", client.UserId))

	manager.DeleteClient(client)
}

// InClient 客户端连接是否存在
func (manager *ClientManager) InClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	_, ok = manager.Clients[client]

	return
}

// GetClients 获取所有连接客户端
func (manager *ClientManager) GetClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)

	manager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value

		return true
	})

	return
}

// ClientsRange 遍历所有客户端, 返回客户端是否存在
func (manager *ClientManager) ClientsRange(f func(client *Client, value bool) (result bool)) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	for key, value := range manager.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}

	return
}

// ClientsCount 获取连接客户端数量
func (manager *ClientManager) ClientsCount() (clientsCount int) {
	clientsCount = len(manager.Clients)

	return
}

// CreateClient 创建客户端连接
func (manager *ClientManager) CreateClient(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// DeleteClient 删除客户端连接
func (manager *ClientManager) DeleteClient(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

// ChanEventStart 管道事件处理
func (manager *ClientManager) ChanEventStart() {
	for {
		select {
		// TODO 建立连接处理
		case client := <-manager.Register:
			manager.EventRegister(client)

		// TODO 断开连接处理
		case client := <-manager.UnRegister:
			manager.EventUnRegister(client)

		}
	}
}
