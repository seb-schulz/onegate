package config

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/go-chi/httplog/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type (
	db struct {
		Dsn   string
		Debug bool
	}

	urlLogin struct {
		Key          []byte
		ExpiresIn    time.Duration
		ValidMethods []string
	}

	session struct {
		Key       string
		ActiveFor time.Duration
	}

	serverKind int

	logger struct {
		Level slog.Level
		File  string
	}

	config struct {
		RelyingParty struct {
			Name    string
			ID      string
			Origins []string
		}
		DB             db
		Session        session
		UrlLogin       urlLogin
		BaseUrl        url.URL
		PrivateAuthKey *ecdsa.PrivateKey
		Server         struct {
			Kind     serverKind
			HttpPort string
			Limit    struct {
				RequestLimit int
				WindowLength time.Duration
			}
		}
		Features struct {
			UserRegistration bool
		}
		Logger logger
	}
)

const (
	ServerKindHttp serverKind = iota
	ServerKindFcgi
)

var (
	// FIXME: Stupid hack to avoid failing tests
	StrictHooks bool
	Config      config
	defaultYaml = []byte(`
relyingParty:
  name: ""
  id: ""
  origins: []
db:
  dsn: ""
  debug: false
session:
  key: ""
  activeFor: 2h
urlLogin:
  key: ""
  expiresIn: 30s
  validMethods: ["HS256", "HS384", "HS512"]
baseUrl: "http://localhost:9000"
server:
  kind: "http"
  httpPort: ""
  limit:
    requestLimit: 50
    windowLength: "1m"
features:
  userRegistration: true
logger:
  level: "info"
  file: ""
privateAuthKey: ""
`)
)

func init() {
	viper.SetConfigName(".onegate.conf")
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(defaultYaml))

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ONEGATE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath("/etc/onegate/")
	viper.AddConfigPath(".")

	viper.ReadInConfig()

	if err := viper.Unmarshal(&Config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		stringToKindHookFunc(),
		stringToLogLevelHookFunc(),
		base64StringToBytesHookFunc(),
		stringToURLHookFunc(),
		stringToEcdsaPrivateKeyHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		base64StringToStringHookFunc(),
	))); err != nil {
		log.Fatalln(err)
	}

	if len(Config.RelyingParty.ID) == 0 {
		Config.RelyingParty.ID = Config.BaseUrl.Hostname()
	}

	if len(Config.RelyingParty.Origins) == 0 {
		xs := []string{Config.BaseUrl.Hostname()}
		if Config.BaseUrl.Port() != "" {
			xs = append(xs, fmt.Sprintf("%s://%s", Config.BaseUrl.Scheme, Config.BaseUrl.Host))
		}

		Config.RelyingParty.Origins = xs
	}

	if Config.Server.HttpPort == "" {
		Config.Server.HttpPort = Config.httpPort()
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

func stringToEcdsaPrivateKeyHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(&ecdsa.PrivateKey{}) {
			return data, nil
		}
		privKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(data.(string)))
		if StrictHooks && err != nil {
			return nil, fmt.Errorf("cannot parse private key: %w", err)
		}
		return privKey, nil
	}
}

func stringToKindHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(ServerKindHttp) {
			return data, nil
		}

		switch data.(string) {
		case "http":
			return ServerKindHttp, nil
		case "fcgi":
			return ServerKindFcgi, nil
		default:
			return nil, fmt.Errorf("cannot dedect kind of server type")
		}
	}
}

func stringToLogLevelHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(slog.LevelInfo) {
			return data, nil
		}

		return httplog.LevelByName(strings.ToUpper(data.(string))), nil
	}
}

func (c config) httpPort() string {
	if c.Server.Kind != ServerKindHttp {
		return ""
	}

	if c.Server.HttpPort != "" {
		return Config.Server.HttpPort
	}

	if c.BaseUrl.Port() == "" {
		return "9000"
	} else {
		return c.BaseUrl.Port()
	}
}
