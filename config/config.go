package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppPort     int `env:"APP_PORT" envDefault:"5000"`
	AppHttpPort int `env:"APP_HTTP_PORT" envDefault:"5001"`

	TONConfigPath      string `env:"TON_CONFIG_PATH,required"`
	TONContractAddress string `env:"TON_CONTRACT_ADDR,required"`
}

func LoadConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("cannot parse initial ENV vars: %v", err)
	}
	return cfg
}
