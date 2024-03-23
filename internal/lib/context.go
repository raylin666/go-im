package lib

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/transport/http"
	"mt/internal/constant/types"
	"reflect"
)

type httpRequest struct{}

// NewContextHeaderAppID 创建应用权限认证上下文
func NewContextHeaderAppID(ctx context.Context, appid types.HeaderAppID) context.Context {
	return context.WithValue(ctx, types.HeaderAppID{}, appid)
}

// GetContextHeaderAppID 获取应用权限认证上下文
func GetContextHeaderAppID(ctx context.Context) (appid types.HeaderAppID, err error) {
	ctxValue := ctx.Value(types.HeaderAppID{})
	if reflect.TypeOf(ctxValue) != reflect.TypeOf(types.HeaderAppID{}) {
		return appid, errors.New("inconsistent data types")
	}

	return ctxValue.(types.HeaderAppID), nil
}

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
