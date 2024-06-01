package lib

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport/http"
	"reflect"
)

type httpRequest struct{}

// NewContextHttpRequest 创建HTTP请求上下文
func NewContextHttpRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, httpRequest{}, req)
}

// GetContextHttpRequest 获取HTTP请求上下文
func GetContextHttpRequest(ctx context.Context) (req *http.Request) {
	ctxValue := ctx.Value(httpRequest{})
	if reflect.TypeOf(ctxValue) != reflect.TypeOf(&http.Request{}) {
		return nil
	}

	return ctxValue.(*http.Request)
}
