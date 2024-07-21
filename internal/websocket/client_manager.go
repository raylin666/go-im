package websocket

import (
	"go.uber.org/zap"
	"mt/internal/lib"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/model"
	"mt/pkg/utils"
	"sync"
	"time"
)

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

// InClient 客户端连接是否存在
func (clientManager *ClientManager) InClient(client *Client) (ok bool) {
	clientManager.ClientsLock.RLock()
	defer clientManager.ClientsLock.RUnlock()

	_, ok = clientManager.Clients[client]

	return
}

// CreateClient 创建客户端连接
func (clientManager *ClientManager) CreateClient(client *Client) {
	clientManager.ClientsLock.Lock()
	defer clientManager.ClientsLock.Unlock()

	clientManager.Clients[client] = true
}

// DeleteClient 删除客户端连接
func (clientManager *ClientManager) DeleteClient(client *Client) {
	if !clientManager.InClient(client) {

		return
	}

	clientManager.ClientsLock.Lock()
	defer clientManager.ClientsLock.Unlock()
	delete(clientManager.Clients, client)
}

// ClientsCount 获取连接客户端数量
func (clientManager *ClientManager) ClientsCount() (clientsCount int) {
	clientsCount = len(clientManager.Clients)

	return
}

// GetClients 获取所有连接客户端
func (clientManager *ClientManager) GetClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)

	clientManager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value

		return true
	})

	return
}

// ClientsRange 遍历所有客户端, 返回客户端是否存在
func (clientManager *ClientManager) ClientsRange(f func(client *Client, value bool) (result bool)) {
	clientManager.ClientsLock.RLock()
	defer clientManager.ClientsLock.RUnlock()

	for key, value := range clientManager.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}

	return
}

// GetClientsCountByAccount 根据账号获取客户端数量
func (clientManager *ClientManager) GetClientsCountByAccount(accountId string) (count int) {
	count = 0
	if !clientManager.HasAccount(accountId) {

		return
	}

	for c := range clientManager.GetClients() {
		if c.Account.ID != accountId {
			continue
		}

		count++
	}

	return
}

// GetClientsByAccount 根据账号获取所有客户端
func (clientManager *ClientManager) GetClientsByAccount(accountId string) (clients []*Client) {
	clients = make([]*Client, 0)

	if !clientManager.HasAccount(accountId) {

		return
	}

	for c := range clientManager.GetClients() {
		if c.Account.ID != accountId {
			continue
		}

		clients = append(clients, c)
	}

	return
}

// HasAccount 账号是否存在
func (clientManager *ClientManager) HasAccount(accountId string) (ok bool) {
	clientManager.AccountsLock.RLock()
	defer clientManager.AccountsLock.RUnlock()

	_, ok = clientManager.Accounts[accountId]

	return
}

// CreateAccount 创建账号 (同个账号多端登录时, 只存储第一个连接的账号在线信息)
func (clientManager *ClientManager) CreateAccount(account *Account) {
	if clientManager.HasAccount(account.ID) {

		return
	}

	clientManager.AccountsLock.Lock()
	defer clientManager.AccountsLock.Unlock()
	clientManager.Accounts[account.ID] = account
}

// DeleteAccount 删除账号
func (clientManager *ClientManager) DeleteAccount(accountId string) (result bool) {
	clientManager.AccountsLock.Lock()
	defer clientManager.AccountsLock.Unlock()

	if _, ok := clientManager.Accounts[accountId]; ok {
		delete(clientManager.Accounts, accountId)
		result = true
	}

	return
}

// GetAccountsCount 获取所有账号数量
func (clientManager *ClientManager) GetAccountsCount() (count int) {
	count = len(clientManager.Accounts)

	return
}

// GetAccountKeys 获取所有账号ID
func (clientManager *ClientManager) GetAccountKeys() (keys []string) {
	clientManager.AccountsLock.RLock()
	defer clientManager.AccountsLock.RUnlock()

	keys = make([]string, 0)
	for key := range clientManager.Accounts {
		keys = append(keys, key)
	}

	return
}

// GetAccounts 获取所有账号
func (clientManager *ClientManager) GetAccounts() (accounts []string) {
	clientManager.AccountsLock.RLock()
	defer clientManager.AccountsLock.RUnlock()

	accounts = make([]string, 0)
	for _, account := range clientManager.Accounts {
		accounts = append(accounts, account.ID)
	}

	return
}

// EventRegister 建立连接处理
func (clientManager *ClientManager) EventRegister(client *Client) {
	var (
		timeNow = time.Now()

		clientId, _ = utils.GetTCPConnFd(client.Conn.NetConn())

		clientIp = utils.ClientIP(lib.GetContextHttpRequest(client.Ctx))
	)

	// TODO 存储账号连接信息
	accountOnline := &model.AccountOnline{
		AccountId:  client.Account.ID,
		LoginTime:  timeNow,
		LoginIp:    clientIp,
		ClientAddr: client.Conn.RemoteAddr().String(),
		ClientId:   int(clientId),
		DeviceId:   "",
		Os:         model.OsWeb,
		System:     "",
	}

	if accountOnlineErr := dbrepo.NewDefaultDbQuery(DbRepo()).AccountOnline.WithContext(client.Ctx).Create(accountOnline); accountOnlineErr != nil {
		Logger(client.Ctx).Error("创建存储账号连接信息错误", zap.Any("account_online", accountOnline), zap.Error(accountOnlineErr))
	} else {
		client.Account.WithOnlineId(accountOnline.ID)
	}

	// TODO 创建连接
	clientManager.CreateClient(client)

	// TODO 创建账号
	clientManager.CreateAccount(client.Account)

	// TODO 创建账号分布式缓存数据
	SetAccountOnline(client.Account.ID, client.Account)

	Logger(client.Ctx).Info("客户端管理器 - 建立连接事件",
		zap.String("address", client.Addr),
		zap.Any("account", client.Account),
		zap.Any("account_online", accountOnline))
}

// EventUnRegister 断开连接处理
func (clientManager *ClientManager) EventUnRegister(client *Client) {
	var (
		timeNow = time.Now()

		clientIp = utils.ClientIP(lib.GetContextHttpRequest(client.Ctx))
	)

	// TODO 更新账号连接信息
	if client.Account.OnlineId > 0 {
		accountOnlineQuery := dbrepo.NewDefaultDbQuery(DbRepo()).AccountOnline
		_, accountOnlineErr := accountOnlineQuery.WithContext(client.Ctx).Where(accountOnlineQuery.ID.Eq(client.Account.OnlineId)).UpdateSimple(accountOnlineQuery.LogoutTime.Value(timeNow), accountOnlineQuery.LogoutIp.Value(clientIp))
		if accountOnlineErr != nil {
			Logger(client.Ctx).Error("更新账号连接信息错误", zap.Any("account", client.Account), zap.Error(accountOnlineErr))
		}
	}

	// TODO 更新账号信息
	accountQuery := dbrepo.NewDefaultDbQuery(DbRepo()).Account
	_, accountErr := accountQuery.WithContext(client.Ctx).Where(accountQuery.AccountId.Eq(client.Account.ID)).UpdateSimple(accountQuery.Status.Value(model.AccountStatusOffline))
	if accountErr != nil {
		Logger(client.Ctx).Error("更新账号信息错误", zap.Any("account", client.Account), zap.Error(accountErr))
	}

	// TODO 删除连接
	clientManager.DeleteClient(client)

	// TODO 删除账号, 判断账号是否只有此连接; 如果是则删除, 否则不删除
	if clientManager.GetClientsCountByAccount(client.Account.ID) <= 1 {
		clientManager.DeleteAccount(client.Account.ID)
	}
	
	// TODO 删除账号分布式缓存数据
	SetAccountOffline(client.Account.ID)

	Logger(client.Ctx).Info("客户端管理器 - 断开连接事件",
		zap.String("address", client.Addr),
		zap.Any("account", client.Account),
		zap.Time("connect_time", client.ConnectTime),
		zap.Time("heartbeat_time", client.HeartbeatTime))
}

// EventBroadcast 广播消息
func (clientManager *ClientManager) EventBroadcast(message []byte) {

}

// ChanEventHandler 管道事件处理
func (clientManager *ClientManager) ChanEventHandler() {
	for {
		select {
		// TODO 建立连接处理
		case client := <-clientManager.Register:
			clientManager.EventRegister(client)

		// TODO 断开连接处理
		case client := <-clientManager.UnRegister:
			clientManager.EventUnRegister(client)

		// TODO 广播消息处理
		case broadcast := <-clientManager.Broadcast:
			clientManager.EventBroadcast(broadcast)
		}
	}
}
