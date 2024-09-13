package initializers

import (
	"github.com/tongineers/tonlib-go-api/internal/app/dependencies"
	"github.com/tongineers/tonlib-go-api/internal/gateways/web/controllers/apiv1"

	// apiv1Status "github.com/tongineers/tonlib-go-api/internal/gateways/web/controllers/apiv1/status"
	"github.com/gin-gonic/gin"
	apiv1Swagger "github.com/tongineers/tonlib-go-api/internal/gateways/web/controllers/apiv1/swagger"
	"github.com/tongineers/tonlib-go-api/internal/gateways/web/controllers/apiv1/transactions"
	"github.com/tongineers/tonlib-go-api/internal/gateways/web/router"
)

// InitializeRouter initializes new gin router
func InitializeRouter(container dependencies.Container) *gin.Engine {
	r := router.NewRouter()

	ctrls := buildControllers(container)

	for i := range ctrls {
		ctrls[i].DefineRoutes(r)
	}

	return r
}

func buildControllers(container dependencies.Container) []apiv1.Controller {
	return []apiv1.Controller{
		// apiv1Status.NewController(container.BuildInfo),
		transactions.NewController(container.Service, container.Logger),
		apiv1Swagger.NewController(),
	}
}
