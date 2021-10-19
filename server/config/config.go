package config

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type FrontendConfig struct {
	PathToDist string `env:"FRINTEND_PATH"`
	IndexPath  string `env:"FRONTEND_INDEX" envDefault:"index.html"`
}

type HttpServerConfig struct {
	Host         string        `env:"SERVER_HOST" envDefault:""`
	Port         int           `env:"PORT" envDefault:"9080"`
	ShutdownWait time.Duration `env:"SHUTDOWN_WAIT" envDefault:"15s"`
}

type Config struct {
	HttpServer HttpServerConfig
	Frontend   FrontendConfig
}

func New() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", *cfg)
	return cfg
}
