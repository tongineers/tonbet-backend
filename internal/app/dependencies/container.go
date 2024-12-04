package dependencies

import (
	"go.uber.org/zap"

	"github.com/tongineers/dice-ton-api/config"
	repository "github.com/tongineers/dice-ton-api/internal/repositories/bets"
	"github.com/tongineers/dice-ton-api/internal/services/fetcher"
	"github.com/tongineers/dice-ton-api/internal/services/listener"
	"github.com/tongineers/dice-ton-api/internal/services/resolver"
	"github.com/tongineers/dice-ton-api/internal/services/smartcont"
)

// Container is a DI container for application
type Container struct {
	Listener     *listener.Service
	Resolver     *resolver.Service
	Fetcher      *fetcher.Service
	DiceContract *smartcont.Service
	Repository   *repository.Repository
	Config       *config.Config
	Logger       *zap.Logger
}
