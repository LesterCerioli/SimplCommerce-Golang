package database

import (
	"github.com/simplcommerce-go/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database interface {
	GetDB() *gorm.DB
	Close() error
	AutoMigrate(models ...interface{}) error
}

type GormDatabase struct {
	db *gorm.DB
}

func NewDatabase(cfg config.DatabaseConfig) (Database, error) {
	db, err := gorm.Open(postgres.Open(config.DSN(cfg)), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	return &GormDatabase{db: db}, nil
}

func (g *GormDatabase) GetDB() *gorm.DB {
	return g.db
}

func (g *GormDatabase) Close() error {
	sqlDB, err := g.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (g *GormDatabase) AutoMigrate(models ...interface{}) error {
	return g.db.AutoMigrate(models...)
}
