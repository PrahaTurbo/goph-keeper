// Package config provides Config with configuration information for the application.
package config

import (
	"log"

	"github.com/caarlos0/env/v10"
)

// BuildVersion is the version of the build, its value is "N/A" by default
var BuildVersion = "N/A"

// BuildDate is the date of the build, its value is "N/A" by default
var BuildDate = "N/A"

// Config holds the configuration information for the application.
type Config struct {
	Host        string `env:"GKEEPER_SERVER_HOST" envDefault:"localhost"`
	Port        string `env:"GKEEPER_SERVER_PORT" envDefault:"8090"`
	SSLCertPath string `env:"GKEEPER_SSL_CERT_PATH" envDefault:"cert/example.crt"`
}

// LoadConfig parses the provided environment variables.
func LoadConfig() *Config {
	c := &Config{}
	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}

	return c
}
