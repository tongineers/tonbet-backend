package providers

import (
	"github.com/tongineers/dice-ton-api/config"
)

func ConfigProvider() *config.Config {
	return config.LoadConfig()
}
