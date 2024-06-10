package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

const (
	// 客户端连接超时时间
	heartbeatExpirationTime = 60 * time.Second
)

// Client 客户端连接
type Client struct {
	Addr          string          // 客户端地址
	Conn          *websocket.Conn // 连接实例对象
	Send          chan []byte     // 待发送的数据
	FirstTime     time.Time       // 首次连接时间
	HeartbeatTime time.Time       // 上次心跳时间
	AccountId     string          // 账号ID
}

func NewClient(accountId string, conn *websocket.Conn) (client *Client) {
	var currentTime = time.Now()
	client = &Client{
		Addr:          conn.RemoteAddr().String(),
		Conn:          conn,
		Send:          make(chan []byte, 100), // 默认预创建容量为100的消息数据包
		FirstTime:     currentTime,
		HeartbeatTime: currentTime,
		AccountId:     accountId,
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
	if !currentTime.Before(c.HeartbeatTime.Add(heartbeatExpirationTime)) {
		timeout = true
	}

	return
}

// connLoggerFields 获取连接日志字段信息
func (c *Client) connLoggerFields() []zap.Field {
	var addr = zap.String("address", c.Addr)
	var accountId = zap.String("account_id", c.AccountId)
	var firstTime = zap.Time("first_time", c.FirstTime)
	var heartbeatTime = zap.Time("heartbeat_time", c.HeartbeatTime)

	return []zap.Field{addr, accountId, firstTime, heartbeatTime}
}

// Read 读取客户端消息
func (c *Client) Read(ctx context.Context) {
	/*var connLoggerFields = c.connLoggerFields()

	for {
		// c.Conn.ReadMessage 该方法会阻塞等待, 直到收到消息才能继续往下执行
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			Logger(ctx).Error("读取客户端消息失败", connLoggerFields, zap.Error(err))

			return
		}

		// 消息处理
		Logger(ctx).Info("读取客户端消息并开始处理", logAddr, zap.String("message", string(message)))
	}*/
}

// Write 写入客户端消息
func (c *Client) Write(ctx context.Context) {
}
