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

func (manager *ClientManager) Logger() *logger.Logger {
	//TODO implement me

	return manager.Tools.Logger()
}

func (manager *ClientManager) CreateClient(account *Account, conn *websocket.Conn) (client *Client) {
	//TODO implement me

	client = NewClient(manager, account, conn)

	return
}

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

}

// eventListenerHandlerToClientUnRegister 断开客户端连接处理
func (manager *ClientManager) eventListenerHandlerToClientUnRegister(client *Client) {

}

// eventListenerHandlerToMessageBroadcast 广播消息处理
func (manager *ClientManager) eventListenerHandlerToMessageBroadcast(message []byte) {

}
