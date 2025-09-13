package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
	"time"
)

type PostgresConfig struct {
	Password string `env:"POSTGRES_PASS"`
	Username string `env:"POSTGRES_USERNAME"`
	Database string `env:"POSTGRES_DATABASE"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
}

func (c PostgresConfig) GetDSNString() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable", c.Username, c.Password, c.Host, c.Port, c.Database)
}

type AppConfig struct {
	LogLevel   string        `env:"LOG_LEVEL"`
	SessionTTL time.Duration `env:"SESSION_TTL"`
}

type ServerConfig struct {
	Host string `env:"SERVER_HOST"`
	Port string `env:"SERVER_PORT"`
}

type Config struct {
	Postgres PostgresConfig
	App      AppConfig
	Server   ServerConfig
}

var (
	once           sync.Once
	configInstance Config
)

func GetConfig() Config {
	once.Do(func() {
		err := cleanenv.ReadConfig(".env", &configInstance)
		if err != nil {
			panic(fmt.Errorf("failed to read env: %w", err))
		}
	})

	return configInstance
}
