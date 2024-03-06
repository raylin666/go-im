package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"mt/pkg/logger"
	"runtime/debug"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 6 * 60
)

// Client 用户连接
type Client struct {
	ctx    context.Context
	logger *logger.Logger

	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	AppId         uint32          // 登录的平台Id app/web/ios
	UserId        string          // 用户Id (用户登录以后才有)
	FirstTime     uint64          // 首次连接事件
	HeartbeatTime uint64          // 用户上次心跳时间
	LoginTime     uint64          // 登录时间 (用户登录以后才有)
}

func NewClient(
	ctx context.Context,
	logger *logger.Logger,
	addr string,
	socket *websocket.Conn,
	firstTime uint64,
) (client *Client) {
	client = &Client{
		ctx:    ctx,
		logger: logger,

		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100), // 默认预创建容量为100的消息数据包
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}

	return
}

// GetKey 读取客户端数据
func (c *Client) GetKey() (key string) {
	key = GetUserKey(c.AppId, c.UserId)

	return
}

// Read 读取客户端消息
func (c *Client) Read() {
	defer func() {
		if r := recover(); r != nil {
			c.logger.UseApp(c.ctx).Error("读取客户端消息 Recover 失败", zap.Stack(string(debug.Stack())), zap.Any("recover", r))
		}
	}()

	defer func() {
		c.logger.UseApp(c.ctx).Info("读取客户端消息并关闭待发送消息区")
		close(c.Send)
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			c.logger.UseApp(c.ctx).Error("读取客户端消息错误", zap.String("address", c.Addr), zap.Error(err))

			return
		}

		// 处理消息
		c.logger.UseApp(c.ctx).Info("读取客户端消息并处理", zap.String("message", string(message)))

		NewProcess(c.ctx, c.logger).HandlerMessage(c, message)
	}
}

// Write 写入客户端消息
func (c *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			c.logger.UseApp(c.ctx).Error("写入客户端消息 Recover 失败", zap.Stack(string(debug.Stack())), zap.Any("recover", r))
		}
	}()

	defer func() {
		ClientManagerInstance.Unregister <- c
		c.Socket.Close()
		c.logger.UseApp(c.ctx).Info("写入客户端消息 (关闭句柄)")
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 发送消息错误并关闭连接
				c.logger.UseApp(c.ctx).Error("写入客户端消息错误并关闭连接", zap.String("address", c.Addr))

				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
