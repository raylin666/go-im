package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

// Client 用户连接
type Client struct {
	Addr          string          // 客户端地址
	Conn          *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	AppId         uint32          // 登录的平台Id app/web/ios
	UserId        string          // 用户Id (用户登录以后才有)
	FirstTime     uint64          // 首次连接事件
	HeartbeatTime uint64          // 用户上次心跳时间
	LoginTime     uint64          // 登录时间 (用户登录以后才有)
}

func NewClient(conn *websocket.Conn) (client *Client) {
	var currentTime = uint64(time.Now().Unix())
	client = &Client{
		Addr:          conn.RemoteAddr().String(),
		Conn:          conn,
		Send:          make(chan []byte, 100), // 默认预创建容量为100的消息数据包
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
			ManagerInstance().Logger().UseWebSocket(ctx).Error("读取客户端消息异常", logAddr, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
		}
	}()

	defer func() {
		ManagerInstance().Logger().UseWebSocket(ctx).Debug("读取客户端消息结束, 关闭待发送的数据包", logAddr)
		close(c.Send)
	}()

	process := NewProcess(c)
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			ManagerInstance().Logger().UseWebSocket(ctx).Error("读取客户端消息失败", logAddr, zap.Error(err))

			return
		}

		// 消息处理
		ManagerInstance().Logger().UseWebSocket(ctx).Info("读取客户端消息并开始处理", logAddr, zap.String("message", string(message)))

		process.HandlerMessage(ctx, message)
	}
}

// Write 写入客户端消息
func (c *Client) Write(ctx context.Context) {
	var logAddr = zap.String("address", c.Addr)

	defer func() {
		if r := recover(); r != nil {
			ManagerInstance().Logger().UseWebSocket(ctx).Error("写入客户端消息异常", logAddr, zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
		}
	}()

	defer func() {
		ManagerInstance().ClientManager().UnRegister <- c
		c.Conn.Close()
		ManagerInstance().Logger().UseWebSocket(ctx).Debug("写入客户端消息结束, 关闭客户端连接", logAddr)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 写入待发送客户端消息错误并关闭连接
				ManagerInstance().Logger().UseWebSocket(ctx).Error("写入待发送客户端消息错误并关闭连接", logAddr)

				return
			}

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
			ManagerInstance().Logger().UseWebSocket(ctx).Error("发送消息异常", zap.String("address", c.Addr), zap.String("message", string(message)), zap.String("stack", string(debug.Stack())), zap.Any("recover", r))
		}
	}()

	c.Send <- message

	return true
}
