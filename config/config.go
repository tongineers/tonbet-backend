package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		AppPort              int    `env:"APP_PORT"          envDefault:"5000"`
		AppHttpPort          int    `env:"APP_HTTP_PORT"     envDefault:"5001"`
		DBHost               string `env:"DB_HOST"           envDefault:"localhost"`
		DBPort               int32  `env:"DB_PORT"           envDefault:"5432"`
		DBMaxIdleConns       int    `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
		DBMaxOpenConns       int    `env:"DB_MAX_OPEN_CONNS" envDefault:"5"`
		DBName               string `env:"DB_NAME,required"`
		DBUser               string `env:"DB_USER,required"`
		DBPassword           string `env:"DB_PASSWORD,required"`
		TONContractAddr      string `env:"TON_CONTRACT_ADDR,required"`
		TONLastTransactionLT uint64 `env:"TON_LAST_TRANSACTION_LT,required"`
		TONConfigPath        string `env:"TON_CONFIG_PATH,required"`
		TONSecretPath        string `env:"TON_SECRET_PATH,required"`
	}
)

func LoadConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("cannot parse initial ENV vars: %v", err)
	}
	return cfg
}
