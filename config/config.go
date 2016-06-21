package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	log "github.com/Sirupsen/logrus"
)

var environments = map[string]string{
	"production": "production.json",
	"staging":    "staging.json",
	"test":       "test.json",
}

type Config struct {
	PrivateKeyPath     string
	PublicKeyPath      string
	JWTExpirationDelta int
}

var Cfg *Config

func init() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		log.WithField("env", env).Warn("Missing GO_ENV env var defaulting to test environment")
		env = "test"
	}
	var err error
	Cfg, err = NewConfig(env)
	if err != nil {
		log.Fatalf("Error loading config %+v", err.Error())
	}
}

func NewConfig(env string) (*Config, error) {
	_, fileName, _, _ := runtime.Caller(1)
	f := path.Join(path.Dir(fileName), environments[env])
	fmt.Println(f)
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	c := Config{}
	err = json.Unmarshal(data, &c)
	fmt.Println(c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
