package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
)

type Env struct {
	env map[string]interface{}
}

func (c *Env) SetToEnv(key string, value interface{}) {
	c.env[key] = value
}

func (c *Env) GetFromEnvString(key string) string {
	v, err := cast.ToStringE(c.mustGetFromEnv(key))
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Env) GetFromEnvInt(key string) int64 {
	v, err := toInt64(c.mustGetFromEnv(key))
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Env) GetFromEnvBool(key string) bool {
	v, err := cast.ToBoolE(c.mustGetFromEnv(key))
	if err != nil {
		panic(err)
	}
	return v
}

func (c *Env) mustGetFromEnv(key string) interface{} {
	v, ok := c.env[key]
	if !ok {
		panic(fmt.Errorf("not found a key: %s", key))
	}
	return v
}

func (c *Env) ClearEnv() {
	c.env = make(map[string]interface{})
}

func toInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case json.Number:
		val, err := v.Int64()
		if err != nil {
			return 0, err
		}
		return val, nil
	default:
		return cast.ToInt64E(v)
	}
}
