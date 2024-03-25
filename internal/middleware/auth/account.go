package auth

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc/metadata"
	"mt/internal/constant/defined"
	"mt/internal/constant/types"
	"mt/internal/lib"
	"mt/internal/repositories/dbrepo"
	"mt/internal/repositories/dbrepo/model"
	"mt/pkg/repositories"
	"strconv"
)

const (
	// AppKey Headers 头权限应用认证参数名称
	AppKey     = "key"
	AppSecret  = "secret"
	AppUserId  = "userid"
	AppUsersig = "usersig"

	// XMdGlobalKeyName Metadata 元数据传递保存的全局应用权限认证参数名称
	XMdGlobalKeyName     = "x-md-global-key"
	XMdGlobalSecretName  = "x-md-global-secret"
	XMdGlobalUserIdName  = "x-md-global-userid"
	XMdGlobalUsersigName = "x-md-global-usersig"
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
			var (
				appKey    int
				appSecret string
				appUserId string
				usersig   string
			)

			if md, ok := metadata.FromIncomingContext(ctx); ok {
				appKey, _ = strconv.Atoi(md.Get(XMdGlobalKeyName)[0])
				appSecret = md.Get(XMdGlobalSecretName)[0]
				appUserId = md.Get(XMdGlobalUserIdName)[0]
				usersig = md.Get(XMdGlobalUsersigName)[0]
			} else if tr, ok := transport.FromServerContext(ctx); ok {
				appKey, _ = strconv.Atoi(tr.RequestHeader().Get(AppKey))
				appSecret = tr.RequestHeader().Get(AppSecret)
				appUserId = tr.RequestHeader().Get(AppUserId)
				usersig = tr.RequestHeader().Get(AppUsersig)
			}

			if appKey <= 0 || appSecret == "" || appUserId == "" || usersig == "" {
				return nil, defined.ErrorNotVisitAuth
			}

			// TODO 验证应用是否存在
			q := dbrepo.NewDefaultDbQuery(repo).App
			m, _ := q.WithContext(ctx).FirstByKeyAndSecret(uint64(appKey), appSecret)
			modelErr := model.AppAvailableByKeyAndSecret(&m)
			if modelErr != nil {
				return nil, defined.ErrorNotVisitAuth
			}

			// TODO 验证用户是否存在
			tableName := model.AccountTableName(m.Key)
			accountQuery := dbrepo.NewDefaultDbQuery(repo).Account.Table(tableName)
			account, accountErr := accountQuery.WithContext(ctx).FirstByUserId(appUserId)
			if accountErr != nil {
				return nil, defined.ErrorAccountNotFound
			}

			// TODO 存储鉴权验证信息到上下文
			ctx = lib.NewContextHeaderAppID(ctx, types.HeaderAppID{
				Key:     m.Key,
				Secret:  m.Secret,
				UserId:  account.UserId,
				Usersig: usersig,
			})

			return handler(ctx, req)
		}
	}
}
