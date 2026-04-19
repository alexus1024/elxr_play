package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port                 int           `required:"true"`
	Debug                bool          `required:"true"`
	LogLevel             string        `required:"true" split_words:"true"`
	MaxCounter           string        `default:"10" split_words:"true"`
	InitCounter          string        `default:"0" split_words:"true"`
	CounterResponseDelay time.Duration `default:"0s" split_words:"true"`
	ShutdownTimeout      time.Duration `default:"10s" split_words:"true"`
	NatsUrl              string        `required:"true" split_words:"true"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("API", &cfg); err != nil {
		return Config{}, fmt.Errorf("process env config: %w", err)
	}

	return cfg, nil
}
