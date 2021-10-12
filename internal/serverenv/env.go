package serverenv

import (
	userDB "github.com/zacscoding/echo-gorm-realworld-app/internal/user/database"
	"gorm.io/gorm"
)

type ServerEnv struct {
	db     *gorm.DB
	userDB userDB.UserDB
}

type Option func(env *ServerEnv)

// NewServerEnv returns a new ServerEnv applied given options
func NewServerEnv(opts ...Option) *ServerEnv {
	env := &ServerEnv{}
	for _, opt := range opts {
		opt(env)
	}
	return env
}

// WithDB sets db to ServerEnv.
func WithDB(db *gorm.DB) Option {
	return func(env *ServerEnv) {
		env.db = db
	}
}

// WithUserDB sets database.UserDB to ServerEnv.
func WithUserDB(userDB userDB.UserDB) Option {
	return func(env *ServerEnv) {
		env.userDB = userDB
	}
}

// GetDB returns a gorm.DB in ServerEnv.
func (se *ServerEnv) GetDB() *gorm.DB {
	return se.db
}

// GetUserDB returns a database.UserDB in ServerEnv.
func (se *ServerEnv) GetUserDB() userDB.UserDB {
	return se.userDB
}

// Close shuts down this server environments.
func (se *ServerEnv) Close() error {
	return nil
}
