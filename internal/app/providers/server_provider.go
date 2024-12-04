package providers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/tongineers/dice-ton-api/config"
	"github.com/tongineers/dice-ton-api/internal/app/factories"
)

func ServerProvider(router *gin.Engine, conf *config.Config, logger *zap.Logger) *factories.Server {
	return factories.ServerFactory(&factories.ServerConfig{
		HttpPort:   conf.AppHttpPort,
		GrpcPort:   conf.AppPort,
		Router:     router,
		EnableGrpc: false,
	}, logger)
}
