package websocket

import (
	"fmt"
	"mt/internal/app"
	"sync"
	"time"
)

// ClientManager 连接管理
type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 // appId+uuid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	Login       chan *login        // 用户登录处理
	Unregister  chan *Client       // 断开连接处理程序
	Broadcast   chan []byte        // 广播消息-向全部成员发送数据
}

// NewClientManager 初始化连接管理
func NewClientManager() *ClientManager {
	return &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 1000),
		Login:      make(chan *login, 1000),
		Unregister: make(chan *Client, 1000),
		Broadcast:  make(chan []byte, 1000),
	}
}

// GetUserKey 获取用户key
func GetUserKey(appId uint32, userId string) (key string) {
	key = fmt.Sprintf("%d_%s", appId, userId)
	return
}

/**************************  Manager  **************************/

// InClient 客户端连接是否存在
func (manager *ClientManager) InClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	// 连接是否存在
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

// ClientsRange 遍历连接客户端
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

// GetClientsLen 获取连接客户端数量
func (manager *ClientManager) GetClientsLen() (clientsLen int) {

	clientsLen = len(manager.Clients)

	return
}

// AddClients 添加客户端
func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// DelClients 删除客户端
func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

// GetUserClient 获取用户的连接
func (manager *ClientManager) GetUserClient(appId uint32, userId string) (client *Client) {

	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	userKey := GetUserKey(appId, userId)
	if value, ok := manager.Users[userKey]; ok {
		client = value
	}

	return
}

// GetUsersLen
func (manager *ClientManager) GetUsersLen() (userLen int) {
	userLen = len(manager.Users)

	return
}

// 添加用户
func (manager *ClientManager) AddUsers(key string, client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	manager.Users[key] = client
}

// 删除用户
func (manager *ClientManager) DelUsers(key string) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	fmt.Println("DelUsers 4")
	if _, ok := manager.Users[key]; ok {
		delete(manager.Users, key)
	}
}

// 获取用户的key
func (manager *ClientManager) GetUserKeys() (userKeys []string) {

	userKeys = make([]string, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for key := range manager.Users {
		userKeys = append(userKeys, key)
	}

	return
}

// 获取用户的key
func (manager *ClientManager) GetUserList() (userList []string) {

	userList = make([]string, 0)

	app.ClientManager.UserLock.RLock()
	defer app.ClientManager.UserLock.RUnlock()

	for _, v := range app.ClientManager.Users {
		userList = append(userList, v.UserId)
		fmt.Println("GetUserList", v.AppId, v.UserId, v.Addr)
	}

	fmt.Println("GetUserList", app.ClientManager.Users)

	return
}

// 获取用户的key
func (manager *ClientManager) GetUserClients() (clients []*Client) {

	clients = make([]*Client, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for _, v := range manager.Users {
		clients = append(clients, v)
	}

	return
}

// 向全部成员(除了自己)发送数据
func (manager *ClientManager) sendAll(message []byte, ignore *Client) {

	clients := manager.GetUserClients()
	for _, conn := range clients {
		if conn != ignore {
			conn.SendMsg(message)
		}
	}
}

// 用户建立连接事件
func (manager *ClientManager) EventRegister(client *Client) {
	manager.AddClients(client)

	fmt.Println("EventRegister 用户建立连接", client.Addr)

	// client.Send <- []byte("连接成功")
}

// 用户登录
func (manager *ClientManager) EventLogin(login *login) {

	client := login.Client
	// 连接存在，在添加
	if manager.InClient(client) {
		userKey := login.GetKey()
		manager.AddUsers(userKey, login.Client)
	}

	fmt.Println("EventLogin 用户登录", client.Addr, login.AppId, login.UserId)

	//orderId := utils.GetOrderIdTime()
	//SendUserMessageAll(login.AppId, login.UserId, orderId, model.MessageCmdEnter, "哈喽~")
}

// 用户断开连接
func (manager *ClientManager) EventUnregister(client *Client) {
	manager.DelClients(client)

	// 删除用户连接
	userKey := GetUserKey(client.AppId, client.UserId)
	manager.DelUsers(userKey)

	// 清除redis登录数据
	//userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	//if err == nil {
	//	userOnline.LogOut()
	//	cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	//}

	// 关闭 chan
	// close(client.Send)

	//fmt.Println("EventUnregister 用户断开连接", client.Addr, client.AppId, client.UserId)

	if client.UserId != "" {
		//orderId := utils.GetOrderIdTime()
		//SendUserMessageAll(client.AppId, client.UserId, orderId, model.MessageCmdExit, "用户已经离开~")
	}
}

// 管道处理程序
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Register:
			// 建立连接事件
			manager.EventRegister(conn)

		case login := <-manager.Login:
			// 用户登录
			manager.EventLogin(login)

		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.EventUnregister(conn)

		case message := <-manager.Broadcast:
			// 广播事件
			clients := manager.GetClients()
			for conn := range clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
				}
			}
		}
	}
}

/**************************  manager info  ***************************************/
// 获取管理者信息
func GetManagerInfo(isDebug string) (managerInfo map[string]interface{}) {
	managerInfo = make(map[string]interface{})

	managerInfo["clientsLen"] = app.ClientManager.GetClientsLen()
	managerInfo["usersLen"] = app.ClientManager.GetUsersLen()
	managerInfo["chanRegisterLen"] = len(app.ClientManager.Register)
	managerInfo["chanLoginLen"] = len(app.ClientManager.Login)
	managerInfo["chanUnregisterLen"] = len(app.ClientManager.Unregister)
	managerInfo["chanBroadcastLen"] = len(app.ClientManager.Broadcast)

	if isDebug == "true" {
		addrList := make([]string, 0)
		app.ClientManager.ClientsRange(func(client *Client, value bool) (result bool) {
			addrList = append(addrList, client.Addr)

			return true
		})

		users := app.ClientManager.GetUserKeys()

		managerInfo["clients"] = addrList
		managerInfo["users"] = users
	}

	return
}

// 获取用户所在的连接
func GetUserClient(appId uint32, userId string) (client *Client) {
	client = app.ClientManager.GetUserClient(appId, userId)

	return
}

// 定时清理超时连接
func ClearTimeoutConnections() {
	currentTime := uint64(time.Now().Unix())

	clients := app.ClientManager.GetClients()
	for client := range clients {
		if client.IsHeartbeatTimeout(currentTime) {
			fmt.Println("心跳时间超时 关闭连接", client.Addr, client.UserId, client.LoginTime, client.HeartbeatTime)

			client.Socket.Close()
		}
	}
}

// 获取全部用户
func GetUserList() (userList []string) {
	fmt.Println("获取全部用户")

	userList = app.ClientManager.GetUserList()

	return
}

// 全员广播
func AllSendMessages(appId uint32, userId string, data string) {
	fmt.Println("全员广播", appId, userId, data)

	ignore := app.ClientManager.GetUserClient(appId, userId)
	app.ClientManager.sendAll([]byte(data), ignore)
}
