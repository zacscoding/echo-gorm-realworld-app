package config

import (
	"encoding/json"
	"github.com/jeremywohl/flatten"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

type Config struct {
	ServerConfig ServerConfig `json:"server"`
	JWTConfig    JWTConfig    `json:"jwt"`
	DBConfig     DBConfig     `json:"db"`
}

// MarshalJSON returns a flat json data with masking values such as db password or jwt.secret config.
func (c *Config) MarshalJSON() ([]byte, error) {
	cfg := struct {
		ServerConfig ServerConfig `json:"server"`
		JWTConfig    JWTConfig    `json:"jwt"`
		DBConfig     DBConfig     `json:"db"`
	}{
		ServerConfig: c.ServerConfig,
		JWTConfig:    c.JWTConfig,
		DBConfig:     c.DBConfig,
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

type ServerConfig struct {
	Port         int    `json:"port"`
	Timeout      string `json:"timeout"`
	ReadTimeout  string `json:"readTimeout"`
	WriteTimeout string `json:"writeTimeout"`
}

type JWTConfig struct {
	Secret         string `json:"secret"`
	SessionTimeout string `json:"sessionTimeout"`
}

type DBConfig struct {
	DataSourceName string `json:"dataSourceName"`
	Migrate        struct {
		Enable bool   `json:"enable"`
		Dir    string `json:"dir"`
	} `json:"migrate"`
	Pool struct {
		MaxOpen     int    `json:"maxOpen"`
		MaxIdle     int    `json:"maxIdle"`
		MaxLifetime string `json:"maxLifetime"`
	} `json:"pool"`
}

// Load load configs in given order.
// 1. defaultConfig
// 2. environment having "REALWORLD_APP_" prefix
// 3. load config file from given configPath
func Load(configPath string) (*Config, error) {
	k := koanf.New(".")

	// load from default config
	err := k.Load(confmap.Provider(defaultConfig, "."), nil)
	if err != nil {
		log.Printf("failed to load default config. err: %v", err)
		return nil, err
	}

	// load from env
	err = k.Load(env.Provider("REALWORLD_APP_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "REALWORLD_APP_")), "_", ".", -1)
	}), nil)
	if err != nil {
		log.Printf("failed to load config from env. err: %v", err)
	}

	// load from config file if exist
	if configPath != "" {
		path, err := filepath.Abs(configPath)
		if err != nil {
			log.Printf("failed to get absoulute config path. configPath:%s, err: %v", configPath, err)
			return nil, err
		}
		log.Printf("load config file from %s", path)
		if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
			log.Printf("failed to load config from file. err: %v", err)
			return nil, err
		}
	}

	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "json", FlatPaths: false}); err != nil {
		log.Printf("failed to unmarshal with conf. err: %v", err)
		return nil, err
	}
	return &cfg, err
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
