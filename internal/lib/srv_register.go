package lib

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"mt/internal/app"
	"mt/internal/repositories/redisrepo"
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
	ip    string
}

func NewSrvRegister(ctx context.Context, repo repositories.DataRepo, tools *app.Tools) *SrvRegister {
	return &SrvRegister{
		ctx:   ctx,
		repo:  repo,
		tools: tools,
		ip:    app.ServerIp,
	}
}

func (srv *SrvRegister) Register() bool {
	if srv.repo.RedisRepo().Count() <= 0 {
		return false
	}

	redisClient := redisrepo.NewDefaultClient(srv.repo.RedisRepo())
	isOk, _ := redisClient.Exists(srv.ctx, cacheImServerRegisterSet).Result()
	if isOk > 0 && redisClient.SIsMember(srv.ctx, cacheImServerRegisterSet, srv.ip).Val() == true {
		// IP 已注册

		return true
	}

	_, err := redisClient.SAdd(srv.ctx, cacheImServerRegisterSet, srv.ip).Result()
	if err != nil {
		srv.tools.Logger().UseApp(srv.ctx).Error(fmt.Sprintf("服务IP `%s` 注册失败", srv.ip), zap.Error(err))

		return false
	}

	srv.tools.Logger().UseApp(srv.ctx).Info(fmt.Sprintf("服务IP `%s` 注册成功", srv.ip))

	return true
}

func (srv *SrvRegister) UnRegister() bool {
	if srv.repo.RedisRepo().Count() <= 0 {
		return true
	}

	redisClient := redisrepo.NewDefaultClient(srv.repo.RedisRepo())
	isOk, _ := redisClient.Exists(srv.ctx, cacheImServerRegisterSet).Result()
	if isOk == 0 || redisClient.SIsMember(srv.ctx, cacheImServerRegisterSet, srv.ip).Val() == false {
		return true
	}

	_, err := redisClient.SRem(srv.ctx, cacheImServerRegisterSet, srv.ip).Result()
	if err != nil {
		srv.tools.Logger().UseApp(srv.ctx).Error(fmt.Sprintf("服务IP `%s` 移除失败", srv.ip), zap.Error(err))

		return false
	}

	srv.tools.Logger().UseApp(srv.ctx).Info(fmt.Sprintf("服务IP `%s` 移除成功", srv.ip))

	return true
}

func (srv *SrvRegister) ClientAddress() (addrs []string) {
	if srv.repo.RedisRepo().Count() <= 0 {
		return
	}

	redisClient := redisrepo.NewDefaultClient(srv.repo.RedisRepo())
	addrs, _ = redisClient.SMembers(srv.ctx, cacheImServerRegisterSet).Result()
	return
}
