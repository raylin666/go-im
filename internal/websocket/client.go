package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"time"
)

const (
	// 客户端连接超时时间
	heartbeatExpirationTime = 60
)

// Client 客户端连接
type Client struct {
	Addr          string          // 客户端地址
	Conn          *websocket.Conn // 连接实例对象
	Send          chan []byte     // 待发送的数据
	FirstTime     uint64          // 首次连接时间
	HeartbeatTime uint64          // 上次心跳时间
	AccountId     string          // 账号ID
}

func NewClient(accountId string, conn *websocket.Conn) (client *Client) {
	var currentTime = uint64(time.Now().Unix())
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

func (c *Client) Read(ctx context.Context) {

}

func (c *Client) Write(ctx context.Context) {

}
