package config

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type FrontendConfig struct {
	PathToDist string `env:"FRONTEND_PATH"`
	IndexPath  string `env:"FRONTEND_INDEX" envDefault:"index.html"`
}

type HttpServerConfig struct {
	Host string `env:"SERVER_HOST" envDefault:""`
	Port int    `env:"SERVER_PORT" envDefault:"80"`

	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"15s"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"15s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"60s"`
	ShutdownWait time.Duration `env:"SERVER_SHUTDOWN_WAIT" envDefault:"15s"`
}

type DatabaseConfig struct {
	Url string `env:"DATABASE_URL,required"`
}

type Config struct {
	IsDebug    bool `env:"DEBUG" envDefault:"false"`
	HttpServer HttpServerConfig
	Frontend   FrontendConfig
	Database   DatabaseConfig
}

func New() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", *cfg)
	return cfg
}
