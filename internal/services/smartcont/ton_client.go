package smartcont

import (
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"

	"github.com/tongineers/dice-ton-api/config"
)

func NewTonAPIClient(conf *config.Config) (*ton.APIClient, error) {
	client := liteclient.NewConnectionPool()
	err := client.AddConnectionsFromConfigFile(conf.TONConfigPath)
	if err != nil {
		return nil, err
	}

	api := ton.NewAPIClient(client).WithRetry()
	return api.(*ton.APIClient), nil
}
