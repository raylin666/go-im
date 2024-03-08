package websocket

import (
	"github.com/gorilla/websocket"
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
func (c *Client) Read() {

}

// Write 写入客户端消息
func (c *Client) Write() {

}
