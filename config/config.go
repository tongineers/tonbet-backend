package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppPort     int `env:"APP_PORT" envDefault:"5000"`
	AppHttpPort int `env:"APP_HTTP_PORT" envDefault:"5001"`

	DBHost         string `env:"DB_HOST" envDefault:"localhost"`
	DBPort         int32  `env:"DB_PORT" envDefault:"5432"`
	DBName         string `env:"DB_NAME,required"`
	DBUser         string `env:"DB_USER,required"`
	DBPassword     string `env:"DB_PASSWORD,required"`
	DBMaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	DBMaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"5"`

	TONConfigPath        string `env:"TON_CONFIG_PATH,required"`
	TONContractAddr      string `env:"TON_CONTRACT_ADDR,required"`
	TONLastTransactionLT uint64 `env:"TON_LAST_TRANSACTION_LT,required"`
	TONOwnerKeyFile      string `env:"TON_OWNER_KEY_FILE,required" envDefault:"owner.pk"`
}

func LoadConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("cannot parse initial ENV vars: %v", err)
	}
	return cfg
}
