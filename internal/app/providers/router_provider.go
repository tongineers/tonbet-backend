package providers

import (
	"github.com/gin-gonic/gin"

	"github.com/tongineers/dice-ton-api/internal/app/dependencies"
	"github.com/tongineers/dice-ton-api/internal/app/factories"
	apiv1Account "github.com/tongineers/dice-ton-api/internal/gateways/web/controllers/apiv1/account"
	apiv1Bets "github.com/tongineers/dice-ton-api/internal/gateways/web/controllers/apiv1/bets"
	"github.com/tongineers/dice-ton-api/internal/gateways/web/controllers/apiv1/swagger"
)

func RouterProvider(container *dependencies.Container) *gin.Engine {
	return factories.RouterFactory(
		apiv1Account.New(container.DiceContract, container.Logger),
		apiv1Bets.New(container.Repository, container.Logger),
		swagger.New(),
	)
}
