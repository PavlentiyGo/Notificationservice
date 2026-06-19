package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RabbitMQURL string `envconfig:"RABBITMQ_URL" required:"true"`
	BotToken    string `envconfig:"BOT_TOKEN" required:"true"`
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
