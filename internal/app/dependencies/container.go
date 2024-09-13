package dependencies

import (
	tonlib "github.com/mercuryoio/tonlib-go/v2"
	"github.com/tongineers/tonlib-go-api/config"
	"github.com/tongineers/tonlib-go-api/internal/gateways/web/controllers/apiv1/transactions"
	"github.com/tongineers/tonlib-go-api/internal/services/tonapi"
	"go.uber.org/zap"
)

// Container is a DI container for application
type Container struct {
	Service      *tonapi.Service
	Transactions *transactions.Controller

	Client *tonlib.Client
	Key    *tonlib.InputKey
	Config *config.Config
	Logger *zap.Logger
}
