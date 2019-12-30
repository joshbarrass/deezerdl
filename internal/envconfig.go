package internal

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	DebugMode bool `envconfig:"DEBUG_MODE"`
}

var ErrNoEnvConfig = errors.New("unable to load environment variables")

func GetEnvConfig() (*EnvConfig, error) {
	var config EnvConfig
	if err := envconfig.Process("", &config); err != nil {
		return nil, ErrNoEnvConfig
	}
	return &config, nil
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{}
}
