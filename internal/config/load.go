package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName                string
	DBHost                string
	DBPort                int
	DBPassword            string
	DBUser                string
	APPEnv                string
	APIKey                string
	Port                  int
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
	CarAPIURL             string
	CarStartYear          int
	CarEndYear            int
}

func Load() (*Config, error) {
	err := godotenv.Load(filepath.Join("/app/", ".env"))
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "development" {
		port = 5432
	}

	config := &Config{
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     port,
		APPEnv:     appEnv,
		CarAPIURL:  os.Getenv("CAR_API_URL"),
	}

	portStr := os.Getenv("APP_PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}

		config.Port = port
	}

	carStartYear := os.Getenv("CAR_START_YEAR")
	if carStartYear != "" {
		carStart, err := strconv.Atoi(carStartYear)
		if err != nil {
			return nil, err
		}

		config.CarStartYear = carStart
	}

	carEndYear := os.Getenv("CAR_END_YEAR")
	if carEndYear != "" {
		carEnd, err := strconv.Atoi(carEndYear)
		if err != nil {
			return nil, err
		}

		config.CarEndYear = carEnd
	}

	return config, nil
}
