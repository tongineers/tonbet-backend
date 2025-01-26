//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"

	"github.com/tongineers/tonbet-backend/internal/app/dependencies"
	"github.com/tongineers/tonbet-backend/internal/app/providers"
	repository "github.com/tongineers/tonbet-backend/internal/repositories/bets"
	"github.com/tongineers/tonbet-backend/internal/services/fetcher"
	"github.com/tongineers/tonbet-backend/internal/services/listener"
	"github.com/tongineers/tonbet-backend/internal/services/resolver"
	"github.com/tongineers/tonbet-backend/internal/services/smartcont"
)

func BuildApplication() (*Application, error) {
	wire.Build(
		// Providers
		providers.LogsProvider,
		providers.ConfigProvider,
		providers.RouterProvider,
		providers.ServerProvider,
		providers.StoreProvider,

		// Services
		smartcont.New,
		listener.New,
		resolver.New,
		fetcher.New,

		// Repositories
		repository.New,

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
