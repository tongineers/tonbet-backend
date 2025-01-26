package providers

import (
	"github.com/tongineers/tonbet-backend/config"
	"github.com/tongineers/tonbet-backend/internal/app/factories"
	"gorm.io/gorm"
)

func StoreProvider(config *config.Config) (*gorm.DB, error) {
	return factories.StoreFactory(config)
}
