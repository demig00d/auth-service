package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		HTTP
		Log
		AccessToken
		MongoDB
	}

	// HTTP -.
	HTTP struct {
		Port int `env-required:"true" env:"HTTP_PORT" env-default:"8080"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	// AccessToken -.
	AccessToken struct {
		LifeTime int64  `env-required:"true" env:"ACCESS_TOKEN_LIFETIME"`
		Secret   string `env-required:"true" env:"ACCESS_TOKEN_SECRET"`
	}

	// MongoDB -.
	MongoDB struct {
		Host     string `env-required:"true" env:"MONGO_HOST" env-default:"localhost"`
		Port     int    `env-required:"true" env:"MONGO_PORT" env-default:"27017"`
		User     string `env-required:"true" env:"MONGO_INITDB_ROOT_USERNAME"`
		Password string `env-required:"true" env:"MONGO_INITDB_ROOT_PASSWORD"`
	}
)

// NewConfig returns app config.
func NewConfig() (Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	return *cfg, err
}
