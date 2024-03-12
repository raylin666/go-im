package websocket

import (
	"context"
)

func EventPing(ctx context.Context, client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

}
