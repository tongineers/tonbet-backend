package initializers

import "github.com/tongineers/tonlib-go-api/config"

func InitializeConfig() *config.Config {
	return config.LoadConfig()
}
