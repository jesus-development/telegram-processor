package config

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

const ENV_FILE = ".env"

type Config struct {
	DB     DBConfig
	Openai OpenaiConfig
}

func LoadConfig() *Config {
	c := &Config{}
	c.LoadFromEnv()
	return c
}

func (c *Config) LoadFromEnv() {
	err := godotenv.Load(ENV_FILE)
	if err != nil {
		slog.Info("Error loading .env file")
	}

	c.DB = DBConfig{
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		Host:       os.Getenv("DB_HOST"),
		Port:       os.Getenv("DB_PORT"),
		Additional: os.Getenv("DB_ADDITIONAL"),
	}

	c.Openai = OpenaiConfig{
		ApiKey: os.Getenv("OPENAI_API_KEY"),
	}
}
