package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const DevelopmentMode = false

type DBConfig struct {
	DBClient        string `env:"DB_CLIENT"`
	DBConnectionURI string `env:"DB_CONNECTION_URI"`
}

type Config struct {
	DBConfig
	Port          string `env:"PORT" envDefault:"2909"`
	JwtSigningKey string `env:"JWT_SIGNING_KEY" envDefault:"CSwS88WnQjKGBAEI"`
}

func LoadEnv() (cfg Config, err error) {
	if DevelopmentMode {
		return LoadEnvFromFile()
	}
	err = env.Parse(&cfg)
	if err != nil {
		log.Printf("failed to config load from ENV: %s", err)
		return
	}
	log.Printf("config load from ENV: %+v", cfg)
	return
}

func LoadEnvFromFile() (cfg Config, err error) {
	err = godotenv.Load()
	if err != nil {
		log.Printf("failed to config load from ENV: %s", err)
		return
	}
	cfg = Config{
		DBConfig: DBConfig{
			DBClient:        os.Getenv("DB_CLIENT"),
			DBConnectionURI: os.Getenv("DB_CONNECTION_URI"),
		},
		Port:          os.Getenv("PORT"),
		JwtSigningKey: os.Getenv("JWT_SIGNING_KEY"),
	}
	log.Printf("config load from ENV: %+v", cfg)
	return
}
