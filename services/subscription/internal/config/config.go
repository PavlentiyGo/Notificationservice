package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	QueryTimeout time.Duration `envconfig:"DB_QUERY_TIMEOUT" default:"5s"`
}

func NewConfig() (Config, error) {

	var config Config

	if err := envconfig.Process("", &config); err != nil {
		return Config{}, fmt.Errorf("failed to procces config: %w", err)
	}
	return config, nil

}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return config

}
