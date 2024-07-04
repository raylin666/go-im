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

// CreateClient 创建客户端连接
func (clientManager *ClientManager) CreateClient(client *Client) {
	clientManager.ClientsLock.Lock()
	defer clientManager.ClientsLock.Unlock()

	clientManager.Clients[client] = true
}

// DeleteClient 删除客户端连接
func (clientManager *ClientManager) DeleteClient(client *Client) {
	clientManager.ClientsLock.Lock()
	defer clientManager.ClientsLock.Unlock()

	if _, ok := clientManager.Clients[client]; ok {
		delete(clientManager.Clients, client)
	}
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
