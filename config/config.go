package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppPort     int `env:"APP_PORT" envDefault:"5000"`
	AppHttpPort int `env:"APP_HTTP_PORT" envDefault:"5001"`

	TONLibConfigPath   string `env:"TONLIB_CONFIG_PATH,required"`
	TONContractAddress string `env:"TON_CONTRACT_ADDR,required"`

	PublicKey   string `env:"PUBLIC_KEY,required"`
	SecretKey   string `env:"SECRET_KEY,required"`
	KeyPassword string `env:"KEY_PASSWORD,required"`
}

func LoadConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("cannot parse initial ENV vars: %v", err)
	}
	return cfg
}
