package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"mt/internal/app"
	"mt/internal/constant/types"
	"mt/internal/data"
	"mt/internal/grpc"
	"mt/internal/lib"
	"mt/internal/repositories/dbrepo/model"
	"mt/pkg/logger"
	"mt/pkg/utils"
	"sync"
	"time"
)

var _ WebsocketClientManager = (*ClientManager)(nil)

type WebsocketClientManager interface {
	// 日志库
	Logger() *logger.Logger

	// 创建客户端连接
	CreateClient(ctx context.Context, account *Account, conn *websocket.Conn) (client *Client)
	// 建立客户端连接处理通道
	ClientRegister(client *Client)
	// 断开客户端连接处理通道
	ClientUnRegister(client *Client)
}

// ClientManager 客户端连接管理
type ClientManager struct {
	DataLogicRepo struct {
		Account data.AccountRepo
	}

	GrpcClient   grpc.GrpcClient
	Tools        *app.Tools
	Clients      map[*Client]bool     // 全部客户端连接资源
	ClientsLock  sync.RWMutex         // 客户端链接读写锁
	Accounts     map[string][]*Client // 全部账号
	AccountsLock sync.RWMutex         // 账号读写锁
	Register     chan *Client         // 建立客户端连接处理通道
	UnRegister   chan *Client         // 断开客户端连接处理通道
	Broadcast    chan []byte          // 广播消息-向全部成员发送数据
}

// NewClientManager 初始化客户端连接管理
func NewClientManager(accountRepo data.AccountRepo, grpcClient grpc.GrpcClient, tools *app.Tools) (manager WebsocketClientManager, cleanup func()) {
	var clientManager = &ClientManager{
		GrpcClient: grpcClient,
		Tools:      tools,
		Clients:    make(map[*Client]bool),
		Accounts:   make(map[string][]*Client),
		Register:   make(chan *Client, 1000),
		UnRegister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}

	clientManager.DataLogicRepo.Account = accountRepo

	// TODO 注册事件监听处理器
	go clientManager.RegisterEventListenerHandler()

	// TODO 定时循环客户端连接, 清理无心跳客户端连接
	go clientManager.CronMonitorClientsHeartbeat()

	cleanup = func() {
		// TODO 关闭所有客户端连接, 设置客户端离线
		for c := range clientManager.getClients() {
			c.Account.WithLogoutState(model.AccountOnlineLoginStateServer)

			// 注意: 此处不能直接调用 clientManager.ClientUnRegister 方法
			//      因为该方法是将客户端连接塞进通道, 异步监听处理退出操作, 这就导致在执行 cleanup 清理时可能会先执行资源关闭再执行通道逻辑处理, 造成无法完成通道逻辑处理(因为资源已被关闭了)
			clientManager.eventListenerHandlerToClientUnRegister(c)
		}
	}

	return clientManager, cleanup
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
func (manager *ClientManager) rangeClients(f func(client *Client, value bool) (ok bool)) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	for key, value := range manager.Clients {
		ok := f(key, value)
		if ok == false {
			return
		}
	}

	return
}

// hasAccount 账号是否存在/在线
func (manager *ClientManager) hasAccount(accountId string) (ok bool) {
	manager.AccountsLock.RLock()
	defer manager.AccountsLock.RUnlock()

	_, ok = manager.Accounts[accountId]

	return
}

// addAccount 添加账号在线客户端 (帐号支持多客户端连接)
func (manager *ClientManager) addAccount(client *Client) {
	if manager.hasClient(client) {

		return
	}

	manager.AccountsLock.Lock()
	defer manager.AccountsLock.Unlock()
	manager.Accounts[client.Account.ID] = append(manager.Accounts[client.Account.ID], client)
}

// getAccount 获取帐号在线客户端 (帐号支持多客户端连接)
func (manager *ClientManager) getAccount(accountId string) []*Client {
	manager.AccountsLock.Lock()
	defer manager.AccountsLock.Unlock()

	if clients, ok := manager.Accounts[accountId]; ok {
		return clients
	}

	return make([]*Client, 0)
}

// countAccount 获取帐号在线客户端数量 (帐号支持多客户端连接)
func (manager *ClientManager) countAccount(accountId string) (countAccount int) {
	if manager.hasAccount(accountId) == false {
		return
	}

	return len(manager.getAccount(accountId))
}

// deleteAccount 删除账号在线客户端 (帐号支持多客户端连接)
func (manager *ClientManager) deleteAccount(client *Client) (ok bool) {
	accountId := client.Account.ID
	if manager.hasAccount(accountId) == false {
		return true
	}

	clients := manager.getAccount(accountId)
	if len(clients) <= 0 {
		manager.AccountsLock.Lock()
		defer manager.AccountsLock.Unlock()
		delete(manager.Accounts, accountId)
		return true
	}

	var newClients []*Client
	for _, c := range clients {
		if c == client {
			continue
		}

		newClients = append(newClients, c)
	}

	manager.AccountsLock.Lock()
	defer manager.AccountsLock.Unlock()
	if len(newClients) > 0 {
		manager.Accounts[accountId] = newClients
	} else {
		delete(manager.Accounts, accountId)
	}

	return true
}

// getAccountIds 获取所有在线帐号ID
func (manager *ClientManager) getAccountIds() (keys []string) {
	manager.AccountsLock.RLock()
	defer manager.AccountsLock.RUnlock()

	keys = make([]string, 0)
	for key := range manager.Accounts {
		keys = append(keys, key)
	}

	return
}

// countAccounts 获取所有在线帐号数量
func (manager *ClientManager) countAccounts() (countAccounts int) {
	return len(manager.Accounts)
}

func (manager *ClientManager) Logger() *logger.Logger {
	//TODO implement me

	return manager.Tools.Logger()
}

// CreateClient 创建客户端连接
func (manager *ClientManager) CreateClient(ctx context.Context, account *Account, conn *websocket.Conn) (client *Client) {
	//TODO implement me

	client = NewClient(ctx, manager, account, conn)

	return
}

// ClientRegister 建立客户端连接处理通道
func (manager *ClientManager) ClientRegister(client *Client) {
	// TODO 监听客户端消息, 不断读出客户端发送的消息数据包
	go client.Read()

	// TODO 监听服务端消息, 不断写入服务端发送的消息数据包
	go func() {
		// TODO 断开服务端事件处理
		defer func() {
			client.Account.WithLogoutState(model.AccountOnlineLoginStateNormal)

			manager.ClientUnRegister(client)
		}()

		client.Write()
	}()

	manager.Register <- client
}

// ClientUnRegister 断开客户端连接处理通道
func (manager *ClientManager) ClientUnRegister(client *Client) {
	//TODO implement me

	manager.UnRegister <- client
}

// CronMonitorClientsHeartbeat 定时器监听客户端心跳
func (manager *ClientManager) CronMonitorClientsHeartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var timeNow = time.Now()
		for c := range manager.getClients() {
			if !c.IsHeartbeatTimeout(timeNow) {
				continue
			}

			c.Account.WithLogoutState(model.AccountOnlineLoginStateTimeout)

			manager.ClientUnRegister(c)
		}
	}
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
	// 将客户端连接添加至管理器
	manager.addClient(client)

	// 将客户端连接添加至在线帐号
	manager.addAccount(client)
}

// eventListenerHandlerToClientUnRegister 断开客户端连接处理
func (manager *ClientManager) eventListenerHandlerToClientUnRegister(client *Client) {
	// 关闭客户端连接后会触发关闭读取客户端事件, 回调到 <-manager.UnRegister 通道, 避免二次执行该方法。
	// 例如 CronMonitorClientsHeartbeat 方法中判断客户端超时时,调用了一次该方法, 当该方法执行 client.Conn.Close 后, 会触发关闭读取客户端事件进行二次执行
	// 例如服务退出时会触发 cleanup 资源清理, 也是会出现相同的问题
	if manager.hasClient(client) == false {
		return
	}

	client.Conn.Close()

	// 将客户端连接从管理器中删除
	manager.deleteClient(client)

	// 将客户端连接从在线帐号中移除
	manager.deleteAccount(client)

	// TODO 登出帐号
	clientIp := utils.ClientIP(lib.GetContextHttpRequest(client.Ctx))
	manager.DataLogicRepo.Account.Logout(client.Ctx, client.Account.ID, &types.AccountLogoutRequest{
		OnlineId: client.Account.OnlineId,
		ClientIp: &clientIp,
		State:    client.Account.LogoutState,
	})
}

// eventListenerHandlerToMessageBroadcast 广播消息处理
func (manager *ClientManager) eventListenerHandlerToMessageBroadcast(message []byte) {

}
