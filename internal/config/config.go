package config

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strings"
	"time"

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

	config struct {
		RelyingParty struct {
			Name    string
			ID      string
			Origins []string
		}
		DB       db
		HttpPort string
		Session  session
		UrlLogin urlLogin
		BaseUrl  url.URL
	}
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
httpPort: ""
session:
  key: ""
  activeFor: 2h
urlLogin:
  key: ""
  expiresIn: 30s
  validMethods: ["HS256", "HS384", "HS512"]
baseUrl: "http://localhost:9000"
`)
)

func init() {
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(defaultYaml))

	for _, p := range configPaths() {
		viper.AddConfigPath(p)
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ONEGATE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(&Config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		base64StringToBytesHookFunc(),
		stringToURLHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
		base64StringToStringHookFunc(),
	))); err != nil {
		log.Fatalln(err)
	}

	for k, v := range map[string]bool{"relyingParty.name": Config.RelyingParty.Name == "", "db.dsn": Config.DB.Dsn == "", "session.key": Config.Session.Key == "", "urlLogin.key": len(Config.UrlLogin.Key) == 0} {
		if v {
			log.Fatalf("missing value for %v", k)
		}
	}

	if len(Config.RelyingParty.ID) <= 0 {
		Config.RelyingParty.ID = Config.BaseUrl.Hostname()
	}

	if Config.HttpPort == "" && Config.BaseUrl.Port() != "" {
		Config.HttpPort = Config.BaseUrl.Port()
	} else if Config.HttpPort == "" && Config.BaseUrl.Port() == "" {
		Config.HttpPort = "443"
	}

	if len(Config.RelyingParty.Origins) == 0 {
		xs := []string{Config.BaseUrl.Hostname()}
		if Config.BaseUrl.Port() != "" {
			xs = append(xs, fmt.Sprintf("%s://%s", Config.BaseUrl.Scheme, Config.BaseUrl.Host))
		}

		Config.RelyingParty.Origins = xs
	}
}

func configPaths() []string {
	return []string{"/etc/onegate/", "$HOME/.onegate/"}
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

func (c config) MarshalYAML() (interface{}, error) {
	return struct {
		RelyingParty struct {
			Name    string
			ID      string
			Origins []string
		}
		DB       db
		HttpPort string
		Session  session
		UrlLogin urlLogin
		BaseUrl  string
	}{c.RelyingParty, c.DB, c.HttpPort, c.Session, c.UrlLogin, c.BaseUrl.String()}, nil
}
