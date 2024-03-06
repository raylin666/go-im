package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// HTTP 升级 WebSocket 协议的配置
var upgrader = websocket.Upgrader{
	HandshakeTimeout: 5 * time.Second,
	ReadBufferSize:   2048,
	WriteBufferSize:  2048,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		fmt.Println(status, reason)
	},
}

func NewUpgrader(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	var header = make(http.Header)
	header.Add("X-IM-SERVER-NAME", "goim")
	header.Add("X-IM-SERVER-VERSION", "v1")
	return upgrader.Upgrade(w, r, header)
}
