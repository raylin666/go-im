package websocket

import (
	"github.com/gorilla/websocket"
	"mt/internal/app"
	"mt/pkg/logger"
	"sync"
)

var _ ClientManagerInterface = (*ClientManager)(nil)

type ClientManagerInterface interface {
	// 日志库
	Logger() *logger.Logger

	// 创建客户端连接
	CreateClient(account *Account, conn *websocket.Conn) (client *Client)
	// 建立客户端连接处理通道
	ClientRegister(client *Client)
	// 断开客户端连接处理通道
	ClientUnRegister(client *Client)
}

// ClientManager 客户端连接管理
type ClientManager struct {
	Tools        *app.Tools
	Clients      map[*Client]bool    // 全部客户端连接资源
	ClientsLock  sync.RWMutex        // 客户端链接读写锁
	Accounts     map[string]*Account // 全部账号
	AccountsLock sync.RWMutex        // 账号读写锁
	Register     chan *Client        // 建立客户端连接处理通道
	UnRegister   chan *Client        // 断开客户端连接处理通道
	Broadcast    chan []byte         // 广播消息-向全部成员发送数据
}

// NewClientManager 初始化客户端连接管理
func NewClientManager(tools *app.Tools) (manager *ClientManager) {
	manager = &ClientManager{
		Tools:      tools,
		Clients:    make(map[*Client]bool),
		Accounts:   make(map[string]*Account),
		Register:   make(chan *Client, 1000),
		UnRegister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}

	// 注册事件监听处理器
	go manager.RegisterEventListenerHandler()

	return
}

// addClient 将客户端连接加入至管理器
func (manager *ClientManager) addClient(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// hasClient 客户端连接是否存在于管理器
func (manager *ClientManager) hasClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	_, ok = manager.Clients[client]

	return
}

// deleteClient 将客户端连接从管理器中删除
func (manager *ClientManager) deleteClient(client *Client) {
	if !manager.hasClient(client) {

		return
	}

	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	delete(manager.Clients, client)
}

// countClients 获取管理器中所有连接的客户端数量
func (manager *ClientManager) countClients() (countClients int) {
	countClients = len(manager.Clients)

	return
}

// getClients 获取管理器中所有连接的客户端
func (manager *ClientManager) getClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)

	manager.rangeClients(func(client *Client, value bool) (result bool) {
		clients[client] = value

		return true
	})

	return
}

// rangeClients 遍历所有管理器中的客户端连接, 返回客户端连接是否存在
func (manager *ClientManager) rangeClients(f func(client *Client, value bool) (result bool)) {
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

func (manager *ClientManager) Logger() *logger.Logger {
	//TODO implement me

	return manager.Tools.Logger()
}

// CreateClient 创建客户端连接
func (manager *ClientManager) CreateClient(account *Account, conn *websocket.Conn) (client *Client) {
	//TODO implement me

	client = NewClient(manager, account, conn)

	return
}

// ClientRegister 建立客户端连接处理通道
func (manager *ClientManager) ClientRegister(client *Client) {
	//TODO implement me

	// 监听客户端消息, 不断读出客户端发送的消息数据包
	go client.Read()

	// 监听服务端消息, 不断写入服务端发送的消息数据包
	go func() {
		// 解绑客户端连接
		defer manager.ClientUnRegister(client)

		client.Write()
	}()

	manager.Register <- client
}

// ClientUnRegister 断开客户端连接处理通道
func (manager *ClientManager) ClientUnRegister(client *Client) {
	//TODO implement me

	manager.UnRegister <- client
}

// RegisterEventListenerHandler 注册事件监听处理器
func (manager *ClientManager) RegisterEventListenerHandler() {
	for {
		select {
		// TODO 建立客户端连接处理
		case client := <-manager.Register:
			manager.eventListenerHandlerToClientRegister(client)

		// TODO 断开客户端连接处理
		case client := <-manager.UnRegister:
			manager.eventListenerHandlerToClientUnRegister(client)

		// TODO 广播消息处理
		case message := <-manager.Broadcast:
			manager.eventListenerHandlerToMessageBroadcast(message)
		}
	}
}

// eventListenerHandlerToClientRegister 建立客户端连接处理
func (manager *ClientManager) eventListenerHandlerToClientRegister(client *Client) {
	// 将客户端连接加入至管理器
	manager.addClient(client)
}

// eventListenerHandlerToClientUnRegister 断开客户端连接处理
func (manager *ClientManager) eventListenerHandlerToClientUnRegister(client *Client) {
	// 将客户端连接从管理器中删除
	manager.deleteClient(client)
}

// eventListenerHandlerToMessageBroadcast 广播消息处理
func (manager *ClientManager) eventListenerHandlerToMessageBroadcast(message []byte) {

}
