package initializers

import "github.com/tongineers/dice-ton-api/config"

func InitializeConfig() *config.Config {
	return config.LoadConfig()
}
