package jetton

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tongineers/tonbet-backend/internal/gateways/web/controllers/apiv1"
	"github.com/tongineers/tonbet-backend/internal/services/smartcont"
)

var (
	_ apiv1.Controller = (*Controller)(nil)
)

type (
	Controller struct {
		dice   *smartcont.Service
		logger *zap.Logger
	}
)

func New(
	dice *smartcont.Service,
	logger *zap.Logger,
) *Controller {
	return &Controller{
		dice:   dice,
		logger: logger,
	}
}

// GetJetton godoc
// @Summary Get Jetton Data
// @Description get jetton
// @ID get-jetton
// @Accept json
// @Produce json
// @Param addr path string true "Account Address"
// @Success 200
// @Router /api/v1/jetton/{addr} [get]
func (ctrl *Controller) GetJetton(ctx *gin.Context) {
	addr := ctx.Param("addr")
	data, err := ctrl.dice.GetJettonData(addr)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, data)
}

func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	r.GET("/api/v1/jetton/:addr", ctrl.GetJetton)
}
