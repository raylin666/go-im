package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"mt/pkg/logger"
	"runtime/debug"
	"time"
)

const (
	// 客户端连接超时时间
	heartbeatExpirationTime = 60 * time.Second
)

// Client 客户端连接
type Client struct {
	Ctx           context.Context
	Manager       ClientManagerInterface
	Addr          string          // 客户端地址
	Conn          *websocket.Conn // 连接实例对象
	Send          chan []byte     // 待发送的数据
	ConnectTime   time.Time       // 客户端连接时间
	HeartbeatTime time.Time       // 上次心跳时间
	Account       *Account        // 账号信息
}

func NewClient(ctx context.Context, manager ClientManagerInterface, account *Account, conn *websocket.Conn) (client *Client) {
	var currentTime = time.Now()
	client = &Client{
		Ctx:           ctx,
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

func (c *Client) logger() *logger.Logger { return c.Manager.Logger() }

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
	var loggerFields = []zap.Field{
		zap.String("address", c.Addr),
		zap.Any("account", c.Account),
		zap.Time("connect_time", c.ConnectTime),
		zap.Time("heartbeat_time", c.HeartbeatTime),
	}

	defer func() {
		if r := recover(); r != nil {
			loggerFields = append(loggerFields, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
			c.logger().Error("读取客户端消息异常", loggerFields...)
		}
	}()

	defer func() {
		// 关闭接收及待发送消息
		close(c.Send)

		c.logger().Info("读取客户端消息结束, 已关闭数据接收", loggerFields...)
	}()

	for {
		// c.Conn.ReadMessage 该方法会阻塞等待, 直到收到消息才能继续往下执行
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				loggerFields = append(loggerFields, zap.String("close_desc", "socket 连接已被关闭"))
			}

			loggerFields = append(loggerFields, zap.Error(err))
			c.logger().Error("读取客户端消息失败", loggerFields...)

			return
		}

		loggerFields = append(loggerFields, zap.String("message", string(message)))
		c.logger().Info("读取客户端消息成功", loggerFields...)

		// 事件消息处理
		// c.EventMessageHandler(ctx, message, true)
	}
}

// Write 写入客户端消息
func (c *Client) Write() {
	var loggerFields = []zap.Field{
		zap.String("address", c.Addr),
		zap.Any("account", c.Account),
		zap.Time("connect_time", c.ConnectTime),
		zap.Time("heartbeat_time", c.HeartbeatTime),
	}

	defer func() {
		if r := recover(); r != nil {
			loggerFields = append(loggerFields, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
			c.logger().Error("写入客户端消息异常", loggerFields...)
		}
	}()

	defer func() {
		c.logger().Info("写入客户端消息结束, 已关闭客户端连接", loggerFields...)

		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 写入待发送客户端消息错误并关闭连接
				c.logger().Error("写入待发送客户端消息错误, 客户端连接将关闭", loggerFields...)

				return
			}

			// 将消息推送至客户端
			c.Conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
