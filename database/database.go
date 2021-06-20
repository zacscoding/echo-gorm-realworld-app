package database

import (
	"github.com/zacscoding/echo-gorm-realworld-app/config"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// NewDatabase creates a new gorm.DB from given config.
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	var (
		db     *gorm.DB
		err    error
		logger = NewLogger(time.Second, true, zapcore.InfoLevel)
	)

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(cfg.DBConfig.DataSourceName), &gorm.Config{
			Logger: logger,
		})
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.DBConfig.Pool.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.DBConfig.Pool.MaxIdle)
	poolMaxLifetime, err := time.ParseDuration(cfg.DBConfig.Pool.MaxLifetime)
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(poolMaxLifetime)

	if cfg.DBConfig.Migrate.Enable {
		err := migrateDB(cfg.DBConfig.DataSourceName, cfg.DBConfig.Migrate.Dir)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
