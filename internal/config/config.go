package config

import (
	"errors"
	"os"
)

type Config struct {
	DBURL string
	Port  string
	Env   string
}

var AppConfig *Config

func LoadConfig() error {

	config := &Config{}

	config.DBURL = os.Getenv("DBURL")
	if config.DBURL == "" {
		return errors.New("DB URL not provided in environment variables")
	}

	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080" // Default port
	}

	config.Env = os.Getenv("ENV")
	if config.Env == "" {
		config.Env = "DEV" // Default environment
	}

	AppConfig = config
	return nil
}
