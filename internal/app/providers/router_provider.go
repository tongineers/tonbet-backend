package providers

import (
	"github.com/gin-gonic/gin"

	"github.com/tongineers/tonbet-backend/internal/app/dependencies"
	"github.com/tongineers/tonbet-backend/internal/app/factories"
	apiv1Account "github.com/tongineers/tonbet-backend/internal/gateways/web/controllers/apiv1/account"
	apiv1Bets "github.com/tongineers/tonbet-backend/internal/gateways/web/controllers/apiv1/bets"
	"github.com/tongineers/tonbet-backend/internal/gateways/web/controllers/apiv1/swagger"
)

func RouterProvider(container *dependencies.Container) *gin.Engine {
	return factories.RouterFactory(
		apiv1Account.New(container.DiceContract, container.Logger),
		apiv1Bets.New(container.Repository, container.Logger),
		swagger.New(),
	)
}
