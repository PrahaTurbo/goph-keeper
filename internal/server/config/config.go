// Package config provides the configurations for both the server and PostgreSQL database.
package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v10"
	"gopkg.in/yaml.v3"
)

const path = "./server.config.yml"

// LoadConfig loads the server and database configurations from a YAML file
// and environment variables into Config struct.
func LoadConfig() *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}

// Config is a representation of the configuration in the YAML file.
// It holds configurations for the server and PostgreSQL database.
type Config struct {
	Server Server `yaml:"goph-keeper"`
	PG     PG     `yaml:"postgre"`
}

// Server holds the server configurations.
type Server struct {
	Host     string `yaml:"host"`
	CertPath string `yaml:"cert_path"`
	KeyPath  string `yaml:"key_path"`
	Secret   string `env:"GKEEPER_SECRET_KEY" envDefault:"secret_key"`
	Port     int    `yaml:"port"`
}

// PG holds the PostgreSQL database configurations.
type PG struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Port     int    `yaml:"port"`
}
