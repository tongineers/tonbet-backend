package smartcont

import (
	"github.com/xssnick/tonutils-go/liteclient"
)

func NewConnectionPool(configFile string) (*liteclient.ConnectionPool, error) {
	client := liteclient.NewConnectionPool()

	if err := client.AddConnectionsFromConfigFile(configFile); err != nil {
		return nil, err
	}

	return client, nil
}
