package serverenv

import (
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/database"
	userDB "github.com/zacscoding/echo-gorm-realworld-app/internal/user/database"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
)

func SetupWith(cfg *config.Config) (*ServerEnv, error) {
	logger := logging.DefaultLogger()
	logger.Info("Setting up application environments.")

	var opts []Option

	// Setup database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		logger.Errorw("failed to initalize database.", "err", err)
		return nil, err
	}
	opts = append(opts, WithDB(db))

	// Setup userDB
	opts = append(opts, WithUserDB(userDB.NewUserDB(cfg, db)))

	return NewServerEnv(opts...), nil
}