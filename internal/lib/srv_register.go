package lib

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/app"
	"mt/internal/repositories/redisrepo"
	"mt/internal/repositories/redisrepo/action"
	"mt/pkg/repositories"
)

/**
服务注册器
*/

var (
	cacheImServerRegisterSet = "im:server:register"
)

type SrvRegister struct {
	ctx   context.Context
	repo  repositories.DataRepo
	tools *app.Tools
	addr  string
}

func NewSrvRegister(ctx context.Context, addr string, repo repositories.DataRepo, tools *app.Tools) *SrvRegister {
	return &SrvRegister{
		ctx:   ctx,
		repo:  repo,
		tools: tools,
		addr:  addr,
	}
}

func (srv *SrvRegister) Register() bool {
	if srv.repo.RedisRepo().Count() <= 0 {
		return false
	}

	redisClient := redisrepo.NewDefaultClient(srv.repo.RedisRepo())

	lock := action.NewLock(srv.ctx, redisClient, "srv_register")
	if lock.Lock() {
		defer lock.UnLock()

		isOk, _ := redisClient.Exists(srv.ctx, cacheImServerRegisterSet).Result()
		if isOk > 0 && redisClient.SIsMember(srv.ctx, cacheImServerRegisterSet, srv.addr).Val() == true {
			// IP 已注册

			return true
		}

		_, err := redisClient.SAdd(srv.ctx, cacheImServerRegisterSet, srv.addr).Result()
		if err != nil {
			srv.tools.Logger().UseApp(srv.ctx).Error(fmt.Sprintf("服务地址 `%s` 注册失败", srv.addr), zap.Error(err))

			return false
		}

		srv.tools.Logger().UseApp(srv.ctx).Info(fmt.Sprintf("服务地址 `%s` 注册成功", srv.addr))
		return true
	}

	return false
}

func (srv *SrvRegister) UnRegister() bool {
	redisClient := redisrepo.NewDefaultClient(srv.repo.RedisRepo())

	lock := action.NewLock(srv.ctx, redisClient, "srv_unregister")
	if lock.Lock() {
		defer lock.UnLock()

		isOk, _ := redisClient.Exists(srv.ctx, cacheImServerRegisterSet).Result()
		if isOk == 0 || redisClient.SIsMember(srv.ctx, cacheImServerRegisterSet, srv.addr).Val() == false {
			return true
		}

		_, err := redisClient.SRem(srv.ctx, cacheImServerRegisterSet, srv.addr).Result()
		if err != nil {
			srv.tools.Logger().UseApp(srv.ctx).Error(fmt.Sprintf("服务地址 `%s` 移除失败", srv.addr), zap.Error(err))

			return false
		}

		srv.tools.Logger().UseApp(srv.ctx).Info(fmt.Sprintf("服务地址 `%s` 移除成功", srv.addr))

		return true
	}

	return false
}

func (srv *SrvRegister) ClientAddress() (addrs []string) {
	redisClient := redisrepo.NewDefaultClient(srv.repo.RedisRepo())
	addrs, _ = redisClient.SMembers(srv.ctx, cacheImServerRegisterSet).Result()
	return
}
