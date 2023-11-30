package config

import (
	"github.com/joho/godotenv"
	"os"
)

type DBConfig struct {
	DBClient        string
	DBConnectionURI string
}

type Config struct {
	DBConfig      DBConfig
	Port          string
	JwtSigningKey string
}

func LoadEnv() (cfg Config, err error) {
	err = godotenv.Load(".env")
	if err != nil {
		return
	}
	cfg = Config{
		Port:          os.Getenv("PORT"),
		JwtSigningKey: os.Getenv("JWT_SIGNING_KEY"),
		DBConfig: DBConfig{
			DBClient:        os.Getenv("DB_CLIENT"),
			DBConnectionURI: os.Getenv("DB_CONNECTION_URI"),
		},
	}
	return
}
