package config

import (
	"encoding/json"
	"fmt"
	"github.com/jeremywohl/flatten"
	"github.com/knadh/koanf"
	kjson "github.com/knadh/koanf/parsers/json"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	EnvPrefix = "REALWORLD_APP_"
)

type Option func(k *koanf.Koanf) error

type Config struct {
	C             *koanf.Koanf
	LoggingConfig LoggingConfig `json:"logging"`
	ServerConfig  ServerConfig  `json:"server"`
	JWTConfig     JWTConfig     `json:"jwt"`
	DBConfig      DBConfig      `json:"db"`
	CacheConfig   CacheConfig   `json:"cache"`
}

type LoggingConfig struct {
	Level       int    `json:"level"`
	Encoding    string `json:"encoding"`
	Development bool   `json:"development"`
}

type ServerConfig struct {
	Port         int           `json:"port"`
	Timeout      time.Duration `json:"timeout"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
	Docs         struct {
		Enabled bool   `json:"enabled"`
		Path    string `json:"path"`
	} `json:"docs"`
}

type JWTConfig struct {
	Secret         string        `json:"secret"`
	SessionTimeout time.Duration `json:"sessionTimeout"`
}

type DBConfig struct {
	DataSourceName string `json:"dataSourceName"`
	Migrate        struct {
		Enable bool   `json:"enable"`
		Dir    string `json:"dir"`
	} `json:"migrate"`
	Pool struct {
		MaxOpen     int           `json:"maxOpen"`
		MaxIdle     int           `json:"maxIdle"`
		MaxLifetime time.Duration `json:"maxLifetime"`
	} `json:"pool"`
}

type CacheConfig struct {
	Enabled     bool          `json:"enabled"`
	Prefix      string        `json:"prefix"`
	Type        string        `json:"type"`
	TTL         time.Duration `json:"ttl"`
	RedisConfig RedisConfig   `json:"redis"`
}

type RedisConfig struct {
	Cluster      bool          `json:"cluster"`
	Endpoints    []string      `json:"endpoints"`
	ReadTimeout  time.Duration `json:"readTimeout"`
	WriteTimeout time.Duration `json:"writeTimeout"`
	DialTimeout  time.Duration `json:"dialTimeout"`
	PoolSize     int           `json:"poolSize"`
	PoolTimeout  time.Duration `json:"poolTimeout"`
	MaxConnAge   time.Duration `json:"maxConnAge"`
	IdleTimeout  time.Duration `json:"idleTimeout"`
}

// Load loads configs in given order.
// 1. defaultConfig
// 2. environment having "REALWORLD_APP_" prefix
// 3. load config file from given configPath
func Load(configPath string) (*Config, error) {
	opts := []Option{WithConfigEnv(EnvPrefix)}
	if configPath != "" {
		opts = append(opts, WithConfigFile(configPath))
	}
	return LoadWithOptions(opts...)
}

// LoadWithOptions loads configs with given options.
func LoadWithOptions(opts ...Option) (*Config, error) {
	k := koanf.New(".")
	// default configs.
	if err := k.Load(confmap.Provider(defaultConfig, "."), nil); err != nil {
		return nil, err
	}
	// load from options
	for _, opt := range opts {
		if err := opt(k); err != nil {
			return nil, err
		}
	}
	conf := Config{C: k}
	if err := k.UnmarshalWithConf("", &conf, koanf.UnmarshalConf{Tag: "json", FlatPaths: false}); err != nil {
		return nil, err
	}
	conf.C = k
	return &conf, nil
}

// WithConfigMap overwrites given configMap which "." separated keys to koanf.Koanf.
func WithConfigMap(configMap map[string]interface{}) Option {
	return func(k *koanf.Koanf) error {
		return k.Load(confmap.Provider(configMap, "."), nil)
	}
}

// WithConfigFile overwrites configs read from configFile to koanf.Koanf.
// Currently, supports "yaml" and "json" format.
func WithConfigFile(configFile string) Option {
	return func(k *koanf.Koanf) error {
		path, err := filepath.Abs(configFile)
		if err != nil {
			return err
		}
		var (
			parser koanf.Parser
			ext    = filepath.Ext(path)
		)
		switch ext {
		case ".yaml", ".yml":
			parser = kyaml.Parser()
		case ".json":
			parser = kjson.Parser()
		default:
			return fmt.Errorf("not supported config file extension: %s. full path: %s", ext, configFile)
		}
		return k.Load(file.Provider(path), parser)
	}
}

// WithConfigEnv overwrites configs read from environments to koanf.Koanf.
// Forexample, The env value "{prefix}SERVER_PORT=8000" will overwrite ServerConfig.Port value.
func WithConfigEnv(prefix string) Option {
	return func(k *koanf.Koanf) error {
		return k.Load(env.ProviderWithValue(prefix, ".", func(key string, value string) (string, interface{}) {
			// trim prefix and to lowercase
			key = strings.ToLower(strings.TrimPrefix(key, prefix))
			// replace "_" to "."
			key = strings.Replace(key, "_", ".", -1)
			// if value is array type, then split with "," separator.
			switch k.Get(key).(type) {
			case []interface{}, []string:
				return key, strings.Split(value, ",")
			}
			return key, value
		}), nil)
	}
}

// MarshalJSON returns a flat json data with masking values such as db password or jwt.secret config.
func (c *Config) MarshalJSON() ([]byte, error) {
	cfg := struct {
		ServerConfig ServerConfig `json:"server"`
		JWTConfig    JWTConfig    `json:"jwt"`
		DBConfig     DBConfig     `json:"db"`
		CacheConfig  CacheConfig  `json:"cache"`
	}{
		ServerConfig: c.ServerConfig,
		JWTConfig:    c.JWTConfig,
		DBConfig:     c.DBConfig,
		CacheConfig:  c.CacheConfig,
	}
	data, err := json.Marshal(&cfg)
	if err != nil {
		return nil, err
	}
	flat, err := flatten.FlattenString(string(data), "", flatten.DotStyle)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal([]byte(flat), &m)
	if err != nil {
		return nil, err
	}

	maskKeys := map[string]struct{}{
		// add keys if u want to mask some properties.
		"jwt.secret": {},
	}

	for key, val := range m {
		if v, ok := val.(string); ok {
			m[key] = maskPassword(v)
		}
		if _, ok := maskKeys[key]; ok {
			m[key] = "****"
		}
	}
	return json.Marshal(&m)
}

func maskPassword(val string) string {
	regex := regexp.MustCompile(`^(?P<protocol>.+?//)?(?P<username>.+?):(?P<password>.+?)@(?P<address>.+)$`)
	if !regex.MatchString(val) {
		return val
	}
	matches := regex.FindStringSubmatch(val)
	for i, v := range regex.SubexpNames() {
		if "password" == v {
			val = strings.ReplaceAll(val, matches[i], "****")
		}
	}
	return val
}
