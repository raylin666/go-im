package auth

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
)

// NewAccountAuthServer 账号验证 Server 中间件
func NewAccountAuthServer() func(handler middleware.Handler) middleware.Handler {
	return selector.Server(
		// 账号验证 权限验证
		AccountMiddlewareHandler(),
	).Match(func(ctx context.Context, operation string) bool {
		// 路由白名单过滤 | 返回true表示需要处理权限验证, 返回false表示不需要处理权限验证
		fmt.Println(operation)
		return false
	}).Build()
}

// AccountMiddlewareHandler 账号验证 中间件处理器
func AccountMiddlewareHandler() func(handler middleware.Handler) middleware.Handler {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				defer func() {}()

				fmt.Println(tr, req)
			}

			return handler(ctx, req)
		}
	}
}
