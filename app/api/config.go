package main

type Config struct {
	Port     int    `env:"PORT" envDefault:"8080"`
	Debug    bool   `env:"DEBUG" envDefault:"false"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}
