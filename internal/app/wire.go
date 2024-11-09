//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/tongineers/dice-ton-api/config"
	"github.com/tongineers/dice-ton-api/internal/app/dependencies"
	"github.com/tongineers/dice-ton-api/internal/app/initializers"
	"github.com/tongineers/dice-ton-api/internal/services/tonapi"
)

func BuildApplication() (*Application, error) {
	wire.Build(
		config.LoadConfig,

		tonapi.New,
		initializers.InitializeServer,
		initializers.InitializeTonClient,
		initializers.InitializeTonClientOpts,
		initializers.InitializeLogs,

		wire.Struct(new(dependencies.Container), "*"),
		wire.Struct(new(Application), "*"),
	)

	return &Application{}, nil
}
