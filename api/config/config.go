package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string `required:"true"`
	DBPort     string `required:"true"`
	DBUser     string `required:"true"`
	DBPassword string `required:"true"`
	DBName     string `required:"true"`
	APIPort    string `required:"true"`
}

// LoadConfig from .env file in parent directory
func LoadConfig() (*Config, error) {

	log.Printf("loading config from .env file")
	// load .env file from parent directory (..)
	envPath := filepath.Join("..", ".env")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
		return nil, fmt.Errorf("could not load config: %w", err)
	}

	config := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		APIPort:    os.Getenv("API_PORT"),
	}

	// validate config struct
	v := validator.New()

	if err := v.Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	log.Printf("config loaded successfully")
	log.Printf("db: %s@%s:%s/%s", config.DBUser, config.DBHost, config.DBPort, config.DBName)
	log.Printf("APIport: %s", config.APIPort)

	return config, nil
}

// GetDSN returns the postgresql connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
	)
}
