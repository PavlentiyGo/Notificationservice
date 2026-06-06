package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr string `envconfig:"SUBSCRIPTION_ADDR" required:"true"`

	QueryTimeout time.Duration `envconfig:"DB_QUERY_TIMEOUT" default:"5s"`

	DbUser     string `envconfig:"POSTGRES_USER"     required:"true"`
	DbPassword string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	DbHost     string `envconfig:"POSTGRES_HOST"     required:"true"`
	DbPort     string `envconfig:"POSTGRES_PORT"     required:"true"`
	DbName     string `envconfig:"SUBSCRIPTION_DB"   required:"true"`
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
