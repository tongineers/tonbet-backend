package providers

import (
	"github.com/tongineers/tonbet-backend/config"
)

func ConfigProvider() *config.Config {
	return config.LoadConfig()
}
