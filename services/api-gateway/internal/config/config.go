package config

import "time"

type Config struct {
	Addr            string
	GracefulTimeout time.Duration
}
