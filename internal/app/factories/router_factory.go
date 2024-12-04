package factories

import (
	"github.com/gin-gonic/gin"

	"github.com/tongineers/dice-ton-api/internal/gateways/web/controllers/apiv1"
	"github.com/tongineers/dice-ton-api/internal/gateways/web/router"
)

func RouterFactory(ctrls ...apiv1.Controller) *gin.Engine {
	r := router.NewRouter()

	for i := range ctrls {
		ctrls[i].DefineRoutes(r)
	}

	return r
}
