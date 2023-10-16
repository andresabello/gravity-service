package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName                string
	DBHost                string
	DBPassword            string
	DBUser                string
	APIKey                string
	Port                  int
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
}

func Load() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	config := &Config{
		DBName:                os.Getenv("DB_NAME"),
		DBUser:                os.Getenv("DB_USER"),
		DBPassword:            os.Getenv("DB_PASSWORD"),
		DBHost:                os.Getenv("DB_HOST"),
		APIKey:                os.Getenv("API_KEY"),
		TwitterConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		TwitterConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		TwitterAccessToken:    os.Getenv("TWITTER_ACCESS_TOKEN"),
		TwitterAccessSecret:   os.Getenv("TWITTER_ACCESS_SECRET"),
	}

	portStr := os.Getenv("PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err == nil {
			config.Port = port
		}
	}

	return config, nil
}
