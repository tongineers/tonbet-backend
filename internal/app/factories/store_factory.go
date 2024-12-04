package factories

import (
	"fmt"

	"github.com/tongineers/dice-ton-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// type (
// 	Store struct {
// 		db           *gorm.DB
// 		maxIdleConns int
// 		maxOpenConns int
// 	}
// )

const (
	DefaultBatchSize = 1000
)

func StoreFactory(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(getConnString(config)), &gorm.Config{
		CreateBatchSize: DefaultBatchSize,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.DBMaxOpenConns)

	return db, sqlDB.Ping()
}

func getConnString(conf *config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		conf.DBUser,
		conf.DBPassword,
		conf.DBHost,
		conf.DBPort,
		conf.DBName,
	)
}
