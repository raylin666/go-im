package websocket

import "sync"

// ClientManager 客户端连接管理
type ClientManager struct {
	Clients      map[*Client]bool    // 全部客户端连接资源
	ClientsLock  sync.RWMutex        // 客户端链接读写锁
	Accounts     map[string]*Account // 全部账号
	AccountsLock sync.RWMutex        // 账号读写锁
	Register     chan *Client        // 建立连接处理通道
	UnRegister   chan *Client        // 断开连接处理通道
	Broadcast    chan []byte         // 广播消息-向全部成员发送数据
}

// NewClientManager 初始化客户端连接管理
func NewClientManager() *ClientManager {
	return &ClientManager{
		Clients:    make(map[*Client]bool),
		Accounts:   make(map[string]*Account),
		Register:   make(chan *Client, 1000),
		UnRegister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}
}
