package providers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tongineers/tonbet-backend/config"
	"github.com/tongineers/tonbet-backend/internal/app/factories"
)

func ServerProvider(router *gin.Engine, conf *config.Config, logger *zap.Logger) *factories.Server {
	return factories.ServerFactory(&factories.ServerConfig{
		HttpPort:   conf.AppHttpPort,
		GrpcPort:   conf.AppPort,
		Router:     router,
		EnableGrpc: false,
	}, logger)
}
