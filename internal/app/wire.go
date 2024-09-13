//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/tongineers/tonlib-go-api/config"
	"github.com/tongineers/tonlib-go-api/internal/app/dependencies"
	"github.com/tongineers/tonlib-go-api/internal/app/initializers"
	"github.com/tongineers/tonlib-go-api/internal/gateways/web/controllers/apiv1/transactions"
	"github.com/tongineers/tonlib-go-api/internal/services/tonapi"
)

func BuildApplication() (*Application, error) {
	wire.Build(
		config.LoadConfig,

		tonapi.New,
		transactions.NewController,

		initializers.InitializeTonClient,
		initializers.InitializeTonClientOpts,

		initializers.InitializeLogs,
		initializers.InitializeRouter,
		initializers.InitializeHTTPServerConfig,
		initializers.InitializeHTTPServer,

		wire.Bind(new(transactions.TONClient), new(*tonapi.Service)),
		wire.Struct(new(dependencies.Container), "*"),
		wire.Struct(new(Application), "*"),
	)

	return &Application{}, nil
}
