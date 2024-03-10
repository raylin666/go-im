package websocket

import (
	"context"
	"sync"
)

type DisposeFunc func(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

var (
	eventHandlers        = make(map[string]DisposeFunc)
	eventHandlersRWMutex sync.RWMutex
)

// RegisterHandlerEvent 注册处理事件
func RegisterHandlerEvent(key string, value DisposeFunc) {
	eventHandlersRWMutex.Lock()
	defer eventHandlersRWMutex.Unlock()
	eventHandlers[key] = value

	return
}

// getHandlerEvent 获取处理事件
func getHandlerEvent(key string) (value DisposeFunc, ok bool) {
	eventHandlersRWMutex.RLock()
	defer eventHandlersRWMutex.RUnlock()

	value, ok = eventHandlers[key]

	return
}
