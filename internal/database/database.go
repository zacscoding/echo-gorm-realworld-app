package database

import (
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// NewDatabase creates a new gorm.DB from given config.
func NewDatabase(conf *config.Config) (*gorm.DB, error) {
	var (
		db     *gorm.DB
		err    error
		logger = NewLogger(time.Second, true, zapcore.InfoLevel)
	)

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(conf.DBConfig.DataSourceName), &gorm.Config{
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
	sqlDB.SetMaxOpenConns(conf.DBConfig.Pool.MaxOpen)
	sqlDB.SetMaxIdleConns(conf.DBConfig.Pool.MaxIdle)
	sqlDB.SetConnMaxLifetime(conf.DBConfig.Pool.MaxLifetime)

	if conf.DBConfig.Migrate.Enable {
		err := migrateDB(conf.DBConfig.DataSourceName, conf.DBConfig.Migrate.Dir)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
