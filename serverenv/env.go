package serverenv

import "gorm.io/gorm"

type ServerEnv struct {
	db *gorm.DB
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

// Close shuts down this server environments.
func (e *ServerEnv) Close() error {
	return nil
}
