package dependencies

import (
	"go.uber.org/zap"

	"github.com/tongineers/tonbet-backend/config"
	repository "github.com/tongineers/tonbet-backend/internal/repositories/bets"
	"github.com/tongineers/tonbet-backend/internal/services/fetcher"
	"github.com/tongineers/tonbet-backend/internal/services/listener"
	"github.com/tongineers/tonbet-backend/internal/services/resolver"
	"github.com/tongineers/tonbet-backend/internal/services/smartcont"
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
