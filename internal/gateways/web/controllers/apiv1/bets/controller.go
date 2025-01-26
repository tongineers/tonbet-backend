package bets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tongineers/tonbet-backend/internal/gateways/web/controllers/apiv1"
	"github.com/tongineers/tonbet-backend/internal/repositories/bets"
)

var (
	_ apiv1.Controller = (*Controller)(nil)
)

type (
	Controller struct {
		repo   *bets.Repository
		logger *zap.Logger
	}
)

func New(
	repo *bets.Repository,
	logger *zap.Logger,
) *Controller {
	return &Controller{
		repo:   repo,
		logger: logger,
	}
}

// GetBets godoc
// @Summary Get Bets
// @Description get bets
// @ID get-bets
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/bets [get]
func (ctrl *Controller) GetBets(ctx *gin.Context) {
	bets, err := ctrl.repo.Read()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, bets)
}

// GetBetsByPlayerAddress godoc
// @Summary Get Bets By Player Address
// @Description get bets by player address
// @ID get-bets-by-player-address
// @Accept json
// @Produce json
// @Param addr path string true "Account Address"
// @Success 200
// @Router /api/v1/bets/{addr} [get]
func (ctrl *Controller) GetBetsByPlayerAddress(ctx *gin.Context) {
	addr := ctx.Param("addr")
	bets, err := ctrl.repo.ReadByPlayerAddress(addr)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, bets)
}

func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	r.GET("/api/v1/bets", ctrl.GetBets)
	r.GET("/api/v1/bets/:addr", ctrl.GetBetsByPlayerAddress)
}
