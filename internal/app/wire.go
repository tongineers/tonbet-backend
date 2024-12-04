//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"

	"github.com/tongineers/dice-ton-api/internal/app/dependencies"
	"github.com/tongineers/dice-ton-api/internal/app/providers"
	repository "github.com/tongineers/dice-ton-api/internal/repositories/bets"
	"github.com/tongineers/dice-ton-api/internal/services/fetcher"
	"github.com/tongineers/dice-ton-api/internal/services/listener"
	"github.com/tongineers/dice-ton-api/internal/services/resolver"
	"github.com/tongineers/dice-ton-api/internal/services/smartcont"
)

func BuildApplication() (*Application, error) {
	wire.Build(

		// Providers
		providers.LogsProvider,
		providers.ConfigProvider,
		providers.RouterProvider,
		providers.ServerProvider,
		providers.StoreProvider,

		// Repositories
		repository.New,

		smartcont.NewTonAPIClient,

		// Services
		listener.New,
		resolver.New,
		fetcher.New,
		smartcont.New,

		wire.Bind(new(listener.DiceContract), new(*smartcont.Service)),
		wire.Bind(new(listener.Repository), new(*repository.Repository)),

		wire.Bind(new(resolver.DiceContract), new(*smartcont.Service)),
		wire.Bind(new(resolver.Repository), new(*repository.Repository)),

		wire.Bind(new(fetcher.DiceContract), new(*smartcont.Service)),
		wire.Bind(new(fetcher.Repository), new(*repository.Repository)),

		wire.Struct(new(dependencies.Container), "*"),
		wire.Struct(new(Application), "*"),
	)

	return &Application{}, nil
}
