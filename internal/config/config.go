package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Address string
}

type PGConfig struct {
	Conn     string
	JdbcUrl  string
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type Config struct {
	Server ServerConfig
	PG     PGConfig
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		Server: ServerConfig{
			Address: os.Getenv("SERVER_ADDRESS"),
		},
		PG: PGConfig{
			Conn:     os.Getenv("POSTGRES_CONN"),
			JdbcUrl:  os.Getenv("POSTGRES_JDBC_URL"),
			Username: os.Getenv("POSTGRES_USERNAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Database: os.Getenv("POSTGRES_DATABASE"),
		},
	}, nil
}
