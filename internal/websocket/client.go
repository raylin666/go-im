package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 60
)

// Client 用户连接
type Client struct {
	Addr          string          // 客户端地址
	Conn          *websocket.Conn // 连接实例对象
	Send          chan []byte     // 待发送的数据
	AppKey        uint64          // 应用KEY
	FirstTime     uint64          // 首次连接时间
	HeartbeatTime uint64          // 上次心跳时间

	// TODO 用户相关 (登录之后才有)
	UserId        string // 用户ID
	LoginPlatform uint32 // 登录的平台ID (App/Web/iOS) - 暂时未用到
	LoginTime     uint64 // 用户登录时间
	LoginIp       string // 用户登录IP
}

func NewClient(key uint64, conn *websocket.Conn) (client *Client) {
	var currentTime = uint64(time.Now().Unix())
	client = &Client{
		Addr:          conn.RemoteAddr().String(),
		Conn:          conn,
		Send:          make(chan []byte, 100), // 默认预创建容量为100的消息数据包
		AppKey:        key,
		FirstTime:     currentTime,
		HeartbeatTime: currentTime,
	}

	return
}

// Read 读取客户端消息
func (c *Client) Read(ctx context.Context) {
	var logAddr = zap.String("address", c.Addr)

	defer func() {
		if r := recover(); r != nil {
			Logger(ctx).Error("读取客户端消息异常", logAddr, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
		}
	}()

	defer func() {
		Logger(ctx).Debug("读取客户端消息结束, 关闭待发送的数据包", logAddr)
		c.SendClose()
	}()

	var process = NewProcess(c)
	for {
		// c.Conn.ReadMessage 该方法会阻塞等待, 直到收到消息才能继续往下执行
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			Logger(ctx).Error("读取客户端消息失败", logAddr, zap.Error(err))

			return
		}

		// 消息处理
		Logger(ctx).Info("读取客户端消息并开始处理", logAddr, zap.String("message", string(message)))

		process.HandlerMessage(ctx, message)
	}
}

// Write 写入客户端消息
func (c *Client) Write(ctx context.Context) {
	var logAddr = zap.String("address", c.Addr)

	defer func() {
		if r := recover(); r != nil {
			Logger(ctx).Error("写入客户端消息异常", logAddr, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
		}
	}()

	defer func() {
		ManagerInstance().ClientManager().UnRegister <- c
		c.Conn.Close()
		Logger(ctx).Debug("写入客户端消息结束, 关闭客户端连接", logAddr)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 写入待发送客户端消息错误并关闭连接
				Logger(ctx).Error("写入待发送客户端消息错误, 客户端连接将关闭", logAddr)

				return
			}

			// 将消息推送至客户端
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// SendMessage 发送消息
func (c *Client) SendMessage(ctx context.Context, message []byte) bool {
	if c == nil {
		return false
	}

	defer func() {
		if r := recover(); r != nil {
			Logger(ctx).Error("发送消息异常", zap.String("address", c.Addr), zap.String("message", string(message)), zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
		}
	}()

	c.Send <- message

	return true
}

// SendClose 关闭接收及待发送消息
func (c *Client) SendClose() {
	close(c.Send)
}

// Heartbeat 更新连接心跳时间
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime

	return
}

// IsHeartbeatTimeout 判断连接心跳是否超时
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}

	return
}

// AccountLogin 账号登录
func (c *Client) AccountLogin(userId string, currentTime uint64, loginIp string, loginPlatform uint32) bool {
	c.UserId = userId
	c.FirstTime = currentTime
	c.LoginIp = loginIp
	c.LoginPlatform = loginPlatform
	c.Heartbeat(currentTime)

	return true
}

// AccountLogout 账号登出
func (c *Client) AccountLogout() bool {
	c.UserId = ""
	c.FirstTime = 0
	c.LoginIp = ""
	c.LoginPlatform = 0

	return true
}

// AccountOnline 账号是否在线
func (c *Client) AccountOnline() bool {
	if c.UserId != "" {
		return true
	}

	return false
}
