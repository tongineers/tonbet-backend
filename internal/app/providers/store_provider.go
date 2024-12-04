package providers

import (
	"github.com/tongineers/dice-ton-api/config"
	"github.com/tongineers/dice-ton-api/internal/app/factories"
	"gorm.io/gorm"
)

func StoreProvider(config *config.Config) (*gorm.DB, error) {
	return factories.StoreFactory(config)
}
