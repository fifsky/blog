//go:build wireinject
// +build wireinject

package main

import (
	"app/cmd"
	"app/config"
	"app/connect"
	"app/handler"
	"app/middleware"
	"app/provider"
	"app/router"
	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

var providerSet = wire.NewSet(
	connect.ProviderSet,
	cmd.ProviderSet,
	router.ProviderSet,
	middleware.ProviderSet,
	handler.ProviderSet,
	provider.RepoSet,
)

func Initialize(conf *config.Config) cli.Commands {
	panic(wire.Build(providerSet))
}
