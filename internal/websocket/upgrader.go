package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type UpgraderOption func(opt *upgraderOption)

type upgraderOption struct {
	HandshakeTimeout time.Duration
	ReadBufferSize   int
	WriteBufferSize  int
	// 解决跨域问题
	CheckOrigin func(r *http.Request) bool
	// 错误处理
	Error func(w http.ResponseWriter, r *http.Request, status int, reason error)
}

func WithDsn(dsn string) UpgraderOption {
	return func(opt *upgraderOption) {
		opt.Dsn = dsn
	}
}

func NewUpgrader(w http.ResponseWriter, r *http.Request, opts ...UpgraderOption) (*websocket.Conn, error) {
	var header = make(http.Header)
	header.Add("X-IM-SERVER-NAME", "goim")
	header.Add("X-IM-SERVER-VERSION", "1.0")

	var o = new(upgraderOption)
	for _, opt := range opts {
		opt(o)
	}

	if o.HandshakeTimeout == 0 {
		o.HandshakeTimeout = 5 * time.Second
	}

	if o.ReadBufferSize == 0 {
		o.ReadBufferSize = 2048
	}

	if o.WriteBufferSize == 0 {
		o.WriteBufferSize = 2048
	}

	// HTTP 升级 WebSocket 协议的配置
	var upgrader = websocket.Upgrader{
		HandshakeTimeout: o.HandshakeTimeout,
		ReadBufferSize:   o.ReadBufferSize,
		WriteBufferSize:  o.WriteBufferSize,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			fmt.Println(status, reason)
		},
	}

	return upgrader.Upgrade(w, r, header)
}
