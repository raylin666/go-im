// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"mt/config"
	"mt/internal/api"
	"mt/internal/app"
	"mt/internal/biz"
	"mt/internal/data"
	"mt/internal/grpc"
	"mt/internal/server"
	"mt/internal/service"
	"mt/internal/websocket"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(bootstrap *config.Bootstrap, configServer *config.Server, configData *config.Data, tools *app.Tools) (*kratos.App, func(), error) {
	dataRepo, cleanup, err := data.NewData(configData, tools)
	if err != nil {
		return nil, nil, err
	}
	heartbeatRepo := data.NewHeartbeatRepo(dataRepo, tools)
	heartbeatUsecase := biz.NewHeartbeatUsecase(heartbeatRepo, tools)
	heartbeatService := service.NewHeartbeatService(heartbeatUsecase)
	accountRepo := data.NewAccountRepo(dataRepo, tools)
	accountUsecase := biz.NewAccountUsecase(accountRepo, tools)
	accountService := service.NewAccountService(accountUsecase, tools)
	grpcServer := server.NewGRPCServer(configServer, heartbeatService, accountService, tools)
	grpcClient, cleanup2, err := grpc.NewGrpcClient(tools)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	clientManager := websocket.NewClientManager(grpcClient, tools)
	handler := api.NewHandler(bootstrap, tools, dataRepo, grpcClient, clientManager)
	httpServer := server.NewHTTPServer(configServer, heartbeatService, accountService, tools, handler)
	kratosApp := newApp(tools, grpcServer, httpServer)
	return kratosApp, func() {
		cleanup2()
		cleanup()
	}, nil
}
