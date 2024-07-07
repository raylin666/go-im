package utils

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"net"
	"strings"
)

const (
	httpBaseContentType = "application"
)

// HttpContentType returns the content-type with base prefix.
func HttpContentType(subtype string) string {
	return strings.Join([]string{httpBaseContentType, subtype}, "/")
}

// HttpContentSubtype returns the content-subtype for the given content-type.  The
// given content-type must be a valid content-type that starts with
// but no content-subtype will be returned.
// according rfc7231.
// contentType is assumed to be lowercase already.
func HttpContentSubtype(contentType string) string {
	left := strings.Index(contentType, "/")
	if left == -1 {
		return ""
	}
	right := strings.Index(contentType, ";")
	if right == -1 {
		right = len(contentType)
	}
	if right < left {
		return ""
	}
	return contentType[left+1 : right]
}

// LocalIP 获取当前服务器IP地址
func LocalIP() (ip string) {
	// net.InterfaceAddrs 也可以实现, 但是多网卡时不推荐, 故采用 UDP 方法获取服务器IP地址
	// UDP 不需要关注是否送达, 只需要对应的 ip 和 port 正确, 即可获取到 IP 地址。
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return
	}

	addr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(addr.String(), ":")[0]
	return
}

// ClientIP 获取客户端IP地址
func ClientIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			return ip
		}
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return "127.0.0.1"
}

// GetTCPConnFd 获取TCP连接FD
func GetTCPConnFd(conn net.Conn) (uintptr, error) {
	file, err := conn.(*net.TCPConn).File()
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Fd(), nil
}
