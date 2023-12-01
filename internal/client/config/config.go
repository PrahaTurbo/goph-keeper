package config

import (
	"log"

	"github.com/caarlos0/env/v10"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
)

type Config struct {
	Host        string `env:"GKEEPER_SERVER_HOST" envDefault:"localhost"`
	Port        string `env:"GKEEPER_SERVER_PORT" envDefault:"8090"`
	SSLCertPath string `env:"GKEEPER_SSL_CERT_PATH" envDefault:"cert/example.crt"`
}

func LoadConfig() *Config {
	c := &Config{}
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}

	return c
}
