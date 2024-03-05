package websocket

// 用户登录
type login struct {
	AppId  uint32
	UserId string
	Client *Client
}

// GetKey 读取客户端数据
func (l *login) GetKey() (key string) {
	key = GetUserKey(l.AppId, l.UserId)
	return
}
