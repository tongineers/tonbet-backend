package transactions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tongineers/tonlib-go-api/internal/dto"
	"github.com/tongineers/tonlib-go-api/internal/gateways/web/controllers/apiv1"
	"go.uber.org/zap"
)

var (
	_ apiv1.Controller = (*Controller)(nil)
)

// Controller is a controller implementation for receiving transactions from the TON blockchain.
type Controller struct {
	client TONClient
	logger *zap.Logger
}

// NewController creates new transactions controller instance
func NewController(c TONClient, l *zap.Logger) *Controller {
	return &Controller{
		client: c,
		logger: l,
	}
}

// GetTransactions godoc
// @Summary Get Transactions from the TON Blockchain
// @Description get status
// @ID get-transactions
// @Accept json
// @Produce json
// @Param addr path string true	"Address"
// @Param hash path string true	"Hash"
// @Param lt path int true	"Lt"
// @Success 200 {object} ResponseDoc
// @Router /api/v1/transactions [get]
func (ctrl *Controller) GetTransactions(ctx *gin.Context) {
	var (
		addr string
		hash string
		lt   int
	)

	val, ok := ctx.GetQuery("addr")
	if !ok {
		// ctrl.logger.Error("")

		ctx.Status(http.StatusBadRequest)
		return
	}
	addr = val

	val, ok = ctx.GetQuery("hash")
	if !ok {
		// ctrl.logger.Error("")

		ctx.Status(http.StatusBadRequest)
		return
	}
	hash = val

	val, ok = ctx.GetQuery("lt")
	if !ok {
		// ctrl.logger.Error("")

		ctx.Status(http.StatusBadRequest)
		return
	}

	lt, err := strconv.Atoi(val)
	if err != nil {
		// ctrl.logger.Error("")

		ctx.Status(http.StatusBadRequest)
		return
	}

	txns, err := ctrl.client.GetTransactions(&dto.GetTransactions{
		Addr: addr,
		Hash: hash,
		Lt:   lt,
	})

	ctx.JSON(http.StatusOK, txns)
	return
}

// DefineRoutes adds controller routes to the router
func (ctrl *Controller) DefineRoutes(r gin.IRouter) {
	r.GET("/api/v1/transactions", ctrl.GetTransactions)
}
