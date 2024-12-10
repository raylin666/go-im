package websocket

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

const (
	// 客户端连接超时时间
	heartbeatExpirationTime = 60 * time.Second
)

// Client 客户端连接
type Client struct {
	Manager       ClientManagerInterface
	Addr          string          // 客户端地址
	Conn          *websocket.Conn // 连接实例对象
	Send          chan []byte     // 待发送的数据
	ConnectTime   time.Time       // 客户端连接时间
	HeartbeatTime time.Time       // 上次心跳时间
	Account       *Account        // 账号信息
}

func NewClient(manager ClientManagerInterface, account *Account, conn *websocket.Conn) (client *Client) {
	var currentTime = time.Now()
	client = &Client{
		Manager:       manager,
		Addr:          conn.RemoteAddr().String(),
		Conn:          conn,
		Send:          make(chan []byte, 100), // 默认预创建容量为100的消息数据包
		ConnectTime:   currentTime,
		HeartbeatTime: currentTime,
		Account:       account,
	}

	return
}

// Heartbeat 更新连接心跳时间
func (c *Client) Heartbeat(currentTime time.Time) {
	c.HeartbeatTime = currentTime

	return
}

// IsHeartbeatTimeout 判断连接心跳是否超时
func (c *Client) IsHeartbeatTimeout(currentTime time.Time) (timeout bool) {
	if c.HeartbeatTime.Add(heartbeatExpirationTime).Before(currentTime) {
		timeout = true
	}

	return
}

// Read 读取客户端消息
func (c *Client) Read() {
	var logger = c.Manager.Logger()
	var loggerFields = []zap.Field{
		zap.String("address", c.Addr),
		zap.Any("account", c.Account),
		zap.Time("connect_time", c.ConnectTime),
		zap.Time("heartbeat_time", c.HeartbeatTime),
	}

	defer func() {
		if r := recover(); r != nil {
			loggerFields = append(loggerFields, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
			logger.Error("读取客户端消息异常", loggerFields...)
		}
	}()

	defer func() {
		// 关闭接收及待发送消息
		close(c.Send)

		logger.Debug("读取客户端消息结束, 已关闭数据接收", loggerFields...)
	}()

	for {
		// c.Conn.ReadMessage 该方法会阻塞等待, 直到收到消息才能继续往下执行
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			loggerFields = append(loggerFields, zap.Error(err))
			logger.Error("读取客户端消息失败", loggerFields...)

			return
		}

		loggerFields = append(loggerFields, zap.String("message", string(message)))
		logger.Info("读取客户端消息成功", loggerFields...)

		// 事件消息处理
		// c.EventMessageHandler(ctx, message, true)
	}
}

// Write 写入客户端消息
func (c *Client) Write() {
	var logger = c.Manager.Logger()
	var loggerFields = []zap.Field{
		zap.String("address", c.Addr),
		zap.Any("account", c.Account),
		zap.Time("connect_time", c.ConnectTime),
		zap.Time("heartbeat_time", c.HeartbeatTime),
	}

	defer func() {
		if r := recover(); r != nil {
			loggerFields = append(loggerFields, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
			logger.Error("写入客户端消息异常", loggerFields...)
		}
	}()

	defer func() {
		logger.Debug("写入客户端消息结束, 已关闭客户端连接", loggerFields...)
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 写入待发送客户端消息错误并关闭连接
				logger.Error("写入待发送客户端消息错误, 客户端连接将关闭", loggerFields...)

				return
			}

			// 将消息推送至客户端
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}