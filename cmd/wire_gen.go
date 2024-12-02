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
	"mt/internal/server"
	"mt/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(bootstrap *config.Bootstrap, configServer *config.Server, configData *config.Data, tools *app.Tools) (*kratos.App, func(), error) {
	dataRepo := data.NewDataRepo(tools, configData)
	dataData, cleanup, err := data.NewData(tools, dataRepo)
	if err != nil {
		return nil, nil, err
	}
	heartbeatRepo := data.NewHeartbeatRepo(dataData, tools)
	heartbeatUsecase := biz.NewHeartbeatUsecase(heartbeatRepo, tools)
	heartbeatService := service.NewHeartbeatService(heartbeatUsecase)
	grpcServer := server.NewGRPCServer(configServer, heartbeatService, tools)
	handler := api.NewHandler(bootstrap, tools, dataRepo)
	httpServer := server.NewHTTPServer(configServer, heartbeatService, tools, handler)
	kratosApp := newApp(tools, grpcServer, httpServer)
	return kratosApp, func() {
		cleanup()
	}, nil
}
