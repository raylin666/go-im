package auth

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"mt/internal/constant/defined"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/model"
	"mt/pkg/repositories"
	"strconv"
)

var (
	// 需要权限验证的路由(operation)
	routeOperations = []string{
		"/v1.account.Service/Create",
	}
)

// NewAccountAuthServer 账号验证 Server 中间件
func NewAccountAuthServer(repo repositories.DataRepo) func(handler middleware.Handler) middleware.Handler {
	return selector.Server(
		// 账号验证 权限验证
		AccountMiddlewareHandler(repo),
	).Match(func(ctx context.Context, operation string) bool {
		// 路由白名单过滤 | 返回true表示需要处理权限验证, 返回false表示不需要处理权限验证
		for _, r := range routeOperations {
			if r == operation {
				return true
			}
		}

		return false
	}).Build()
}

// AccountMiddlewareHandler 账号验证 中间件处理器
func AccountMiddlewareHandler(repo repositories.DataRepo) func(handler middleware.Handler) middleware.Handler {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				var header = tr.RequestHeader()

				defer func() {}()

				appKey, err := strconv.Atoi(header.Get("key"))
				appSecret := header.Get("secret")
				usersig := header.Get("usersig")
				if err != nil || appKey <= 0 || appSecret == "" || usersig == "" {
					return nil, defined.ErrorNotVisitAuth
				}

				q := dbrepo.NewDefaultDbQuery(repo).App
				m, _ := q.WithContext(ctx).FirstByKeyAndSecret(uint64(appKey), appSecret)
				modelErr := model.AppAvailableByKeyAndSecret(&m)
				if modelErr != nil {
					return nil, defined.ErrorNotVisitAuth
				}
			}

			return handler(ctx, req)
		}
	}
}
