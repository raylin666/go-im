//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"
	"mt/config"
	"mt/internal/api"
	"mt/internal/app"
	"mt/internal/biz"
	"mt/internal/data"
	"mt/internal/server"
	"mt/internal/service"
	"mt/internal/websocket"
	"mt/internal/websocket/event"
)

// wireApp init kratos application.
func wireApp(*config.Server,
	*config.Data,
	*config.App,
	*config.Websocket,
	*app.Tools) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		api.ProviderSet,
		service.ProviderSet,
		websocket.ProviderSet,
		event.ProviderSet,
		newApp))
}
