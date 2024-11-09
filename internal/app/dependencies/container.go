package dependencies

import (
	"go.uber.org/zap"

	"github.com/tongineers/dice-ton-api/config"
	"github.com/tongineers/dice-ton-api/internal/services/tonapi"
)

// Container is a DI container for application
type Container struct {
	Service *tonapi.Service
	Config  *config.Config
	Logger  *zap.Logger
}
