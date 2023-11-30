package config

import (
	"github.com/caarlos0/env/v10"
	"log"
)

type DBConfig struct {
	DBClient        string `env:"DB_CLIENT"`
	DBConnectionURI string `env:"DB_CONNECTION_URI"`
}

type Config struct {
	DBConfig
	Port          string `env:"PORT" envDefault:"2909"`
	JwtSigningKey string `env:"JWT_SIGNING_KEY"`
}

func LoadEnv() (cfg Config, err error) {
	err = env.Parse(&cfg)
	if err != nil {
		log.Printf("failed to config load from ENV: %s", err)
		return
	}
	log.Printf("config load from ENV: %+v", cfg)
	return
}
