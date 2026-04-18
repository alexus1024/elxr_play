package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port                 int           `env:"PORT,required"`
	Debug                bool          `env:"DEBUG,required"`
	LogLevel             string        `env:"LOG_LEVEL,required"`
	MaxCounter           string        `env:"MAX_COUNTER,default=10"`
	InitCounter          string        `env:"INIT_COUNTER,default=0"`
	CounterResponseDelay time.Duration `env:"COUNTER_RESPONSE_DELAY,default=0s"`
	ShutdownTimeout      time.Duration `env:"SHUTDOWN_TIMEOUT, default=10s"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("API", &cfg); err != nil {
		return Config{}, fmt.Errorf("process env config: %w", err)
	}

	return cfg, nil
}
