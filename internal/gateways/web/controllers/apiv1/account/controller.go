package account

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

// GetAccount godoc
// @Summary Get Account State
// @Description get account
// @ID get-account
// @Accept json
// @Produce json
// @Param addr path string true "Account Address"
// @Success 200
// @Router /api/v1/account/{addr} [get]
func (ctrl *Controller) GetAccount(ctx *gin.Context) {
	addr := ctx.Param("addr")
	account, err := ctrl.dice.GetAccountState(addr)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	r.GET("/api/v1/account/:addr", ctrl.GetAccount)
}
