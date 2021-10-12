package serverenv

import (
	articleDB "github.com/zacscoding/echo-gorm-realworld-app/internal/article/database"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/cache"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/database"
	userDB "github.com/zacscoding/echo-gorm-realworld-app/internal/user/database"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
)

func SetupWith(conf *config.Config) (*ServerEnv, error) {
	logger := logging.DefaultLogger()
	logger.Info("Setting up application environments.")

	var opts []Option

	// Setup database
	db, err := database.NewDatabase(conf)
	if err != nil {
		logger.Errorw("failed to initalize database.", "err", err)
		return nil, err
	}
	opts = append(opts, WithDB(db))

	// Setup userDB
	udb := userDB.NewUserDB(conf, db)
	if conf.CacheConfig.Enabled {
		redisCli, _, err := cache.NewCache(conf)
		if err != nil {
			logger.Errorw("failed to create a redis client", "err", err)
			return nil, err
		}
		udb = userDB.NewUserCacheDB(conf, redisCli, udb)
	}
	opts = append(opts, WithUserDB(udb))

	// Setup articleDB
	adb := articleDB.NewArticleDB(conf, db)
	// TODO: add cache layer DB.
	opts = append(opts, WithArticleDB(adb))

	return NewServerEnv(opts...), nil
}
