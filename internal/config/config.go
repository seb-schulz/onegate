package config

import (
	"bytes"
	"encoding/base64"
	"net/url"
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
  expiresIn: 1h
  validMethods: ["HS256", "HS384", "HS512"]
relyingParty:
  name: "NOT_CONFIGURED_YET"
  id: "NOT_CONFIGURED_YET"
  origins: []
db:
  dsn: "NOT_CONFIGURED_YET"
  debug: false
httpPort: 9000
session:
  key: "NOT_CONFIGURED_YET"
  activeFor: 2h
urlLogin:
  key: ""
  expiresIn: 30s
  validMethods: ["HS256", "HS384", "HS512"]
baseUrl: "http://localhost:9000"
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
	UrlLogin struct {
		Key          []byte
		ExpiresIn    time.Duration
		ValidMethods []string
	}
	BaseUrl url.URL
}

var Default Config

func init() {
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(defaultYaml))
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ONEGATE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(&Default, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		base64StringToBytesHookFunc(),
		stringToURLHookFunc(),
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

func base64StringToBytesHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf([]byte{}) {
			return data, nil
		}

		if result, err := base64.StdEncoding.DecodeString(data.(string)); err == nil {
			return result, nil
		}

		return data, nil
	}
}

func stringToURLHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(url.URL{}) {
			return data, nil
		}

		if result, err := url.Parse(data.(string)); err == nil {
			return result, nil
		}

		return data, nil
	}
}
