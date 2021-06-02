package logging

import "go.uber.org/zap/zapcore"

var defaultCfg = &Config{
	Encoding:    "console",
	Level:       zapcore.DebugLevel,
	Development: true,
}

type Config struct {
	Encoding    string
	Level       zapcore.Level
	Development bool
}

// SetConfig sets given logging configs for DefaultLogger's logger.
// Must set configs before calling DefaultLogger()
func SetConfig(c *Config) {
	defaultCfg = &Config{
		Encoding:    c.Encoding,
		Level:       c.Level,
		Development: c.Development,
	}
}
