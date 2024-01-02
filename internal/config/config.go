package config

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/go-chi/httplog/v2"
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
		DB       db
		Session  session
		UrlLogin urlLogin
		BaseUrl  url.URL
		Server   struct {
			Kind     serverKind
			HttpPort string
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
features:
  userRegistration: true
logger:
  level: "info"
  file: ""
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

func (x db) MarshalYAML() (interface{}, error) {
	return db{base64.StdEncoding.EncodeToString([]byte(x.Dsn)), x.Debug}, nil
}

func (s session) MarshalYAML() (interface{}, error) {
	return session{base64.StdEncoding.EncodeToString([]byte(s.Key)), s.ActiveFor}, nil
}

func (u urlLogin) MarshalYAML() (interface{}, error) {
	return struct {
		Key          string
		ExpiresIn    time.Duration
		ValidMethods []string `yaml:",flow"`
	}{base64.StdEncoding.EncodeToString(u.Key), u.ExpiresIn, u.ValidMethods}, nil
}

func (k serverKind) MarshalYAML() (interface{}, error) {
	if k == ServerKindHttp {
		return "http", nil
	}
	if k == ServerKindFcgi {
		return "fcgi", nil
	}
	return nil, fmt.Errorf("invalid kind of server")
}

func (c config) MarshalYAML() (interface{}, error) {
	return struct {
		RelyingParty struct {
			Name    string
			ID      string
			Origins []string
		}
		DB       db
		Session  session
		UrlLogin urlLogin
		BaseUrl  string
		Server   struct {
			Kind     serverKind
			HttpPort string
		}
		Features struct {
			UserRegistration bool
		}
		Logger logger
	}{c.RelyingParty, c.DB, c.Session, c.UrlLogin, c.BaseUrl.String(), c.Server, c.Features, c.Logger}, nil
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
