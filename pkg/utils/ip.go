package utils

import "net"

// GetServerIp 获取服务器IP
func GetServerIp() (ip string) {
	manyAddress, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, address := range manyAddress {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}

	return
}
