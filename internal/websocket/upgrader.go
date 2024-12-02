package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	// 默认握手超时
	defaultHandshakeTimeout = 5 * time.Second
	// 默认读缓冲大小
	defaultReadBufferSize = 2048
	// 默认写缓冲大小
	defaultWriteBufferSize = 2048
)

type UpgraderResponseHeader struct {
	Name    string
	Version string
}

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

func WithUpgraderHandshakeTimeout(t time.Duration) UpgraderOption {
	return func(opt *upgraderOption) {
		opt.HandshakeTimeout = t
	}
}

func WithUpgraderReadBufferSize(size int) UpgraderOption {
	return func(opt *upgraderOption) {
		opt.ReadBufferSize = size
	}
}

func WithUpgraderWriteBufferSize(size int) UpgraderOption {
	return func(opt *upgraderOption) {
		opt.WriteBufferSize = size
	}
}

func WithUpgraderCheckOrigin(f func(r *http.Request) bool) UpgraderOption {
	return func(opt *upgraderOption) {
		opt.CheckOrigin = func(r *http.Request) bool {
			return f(r)
		}
	}
}

func WithUpgraderError(f func(w http.ResponseWriter, r *http.Request, status int, reason error)) UpgraderOption {
	return func(opt *upgraderOption) {
		opt.Error = func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			f(w, r, status, reason)
		}
	}
}

func NewUpgrader(w http.ResponseWriter, r *http.Request, upHeader *UpgraderResponseHeader, upOpts ...UpgraderOption) (*websocket.Conn, error) {
	var header = make(http.Header)
	if upHeader != nil {
		if upHeader.Name != "" {
			header.Add("X-Server-Name", upHeader.Name)
		}
		if upHeader.Version != "" {
			header.Add("X-Server-Version", upHeader.Version)
		}
	}

	var o = new(upgraderOption)
	for _, opt := range upOpts {
		opt(o)
	}

	if o.HandshakeTimeout == 0 {
		o.HandshakeTimeout = defaultHandshakeTimeout
	}

	if o.ReadBufferSize == 0 {
		o.ReadBufferSize = defaultReadBufferSize
	}

	if o.WriteBufferSize == 0 {
		o.WriteBufferSize = defaultWriteBufferSize
	}

	// HTTP 升级 WebSocket 协议的配置
	var upgrader = websocket.Upgrader{
		HandshakeTimeout: o.HandshakeTimeout,
		ReadBufferSize:   o.ReadBufferSize,
		WriteBufferSize:  o.WriteBufferSize,
		// 解决跨域问题
		CheckOrigin: o.CheckOrigin,
		// 错误处理
		Error: o.Error,
	}

	return upgrader.Upgrade(w, r, header)
}
