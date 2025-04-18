package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	//nolint: golint //reason: blank import because of swagger docs init
	_ "github.com/tongineers/tonbet-backend/api/web"
	"github.com/tongineers/tonbet-backend/internal/gateways/web/controllers/apiv1"
)

var (
	_ apiv1.Controller = (*Controller)(nil)
)

// Controller implements controller for swagger
type Controller struct {
}

// NewController create new instance for swagger controller
func New() *Controller {
	return &Controller{}
}

// DefineRoutes adds swagger controller routes to the router
func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
