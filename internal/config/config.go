package config

import (
	"bytes"
	"encoding/base64"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var defaultYaml = []byte(`
jwt:
  header: x-jwt-token
  secret: "NOT_CONFIGURED_YET"
  expires_in: 1h
  valid_methods: ["HS256", "HS384", "HS512"]
RelyingParty:
  name: "NOT_CONFIGURED_YET"
  id: "NOT_CONFIGURED_YET"
  origins: []
db:
  dsn: "NOT_CONFIGURED_YET"
  debug: false
httpPort: 9000
session:
  key: "NOT_CONFIGURED_YET"
  active_for: 2h
`)

type Config struct {
	RelyingParty struct {
		Name    string
		ID      string
		Origins []string
	}
	JWT struct {
		Header       string
		Secret       string
		ExpiresIn    time.Duration
		ValidMethods []string
	}
	DB struct {
		Dsn   string
		Debug bool
	}
	HttpPort string
	Session  struct {
		Key       string
		ActiveFor time.Duration
	}
}

var Default Config

func init() {
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(defaultYaml))
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ONEGATE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(&Default, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		base64StringToStringHookFunc(),
	))); err != nil {
		panic(err)
	}
}

func base64StringToStringHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t.Kind() != reflect.String {
			return data, nil
		}

		if result, err := base64.StdEncoding.DecodeString(data.(string)); err == nil {
			return result, nil
		}

		return data, nil
	}
}
