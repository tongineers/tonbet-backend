package initializers

import (
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"go.uber.org/zap"

	"github.com/tongineers/dice-ton-api/config"
	"github.com/tongineers/dice-ton-api/internal/services/tonapi"
)

func InitializeTonClient(conf *config.Config) (*ton.APIClient, error) {
	client := liteclient.NewConnectionPool()
	err := client.AddConnectionsFromConfigFile(conf.TONConfigPath)
	if err != nil {
		return nil, err
	}

	api := ton.NewAPIClient(client).WithRetry()
	return api.(*ton.APIClient), nil
}

func InitializeTonClientOpts(
	сlient *ton.APIClient,
	conf *config.Config,
	logger *zap.Logger,
) []tonapi.Opt {
	return []tonapi.Opt{
		tonapi.WithClient(сlient),
		tonapi.WithConfig(conf),
		tonapi.WithLogger(logger),
	}
}
