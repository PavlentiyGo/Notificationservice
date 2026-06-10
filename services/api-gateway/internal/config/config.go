package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerAddr       string        `envconfig:"SERVER_ADDR"       required:"true"`
	SubscriptionAddr string        `envconfig:"SUBSCRIPTION_ADDR" required:"true"`
	GracefulTimeout  time.Duration `envconfig:"GRACEFUL_TIMEOUT"  default:"5s"`

	BotToken string `envconfig:"BOT_TOKEN"         required:"true"`
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
