package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App       AppConfig      `yaml:"app"`
	HTTP      HTTPConfig     `yaml:"http"`
	Services  ServicesConfig `yaml:"services"`
	JWTSecret string         `env:"JWT_SECRET" env-default:"secret"`
}

type AppConfig struct {
	Name string `yaml:"name" env:"APP_NAME" env-default:"gateway-service"`
	Env  string `yaml:"env" env:"APP_ENV" env-default:"development"`
}

type HTTPConfig struct {
	Port string `yaml:"port" env:"APP_PORT" env-default:"8080"`
}

type ServicesConfig struct {
	Auth      string `yaml:"auth" env:"AUTH_SERVICE_ADDR" env-default:"localhost:9000"`
	Catalog   string `yaml:"catalog" env:"CATALOG_SERVICE_ADDR" env-default:"localhost:9001"`
	Orders    string `yaml:"orders" env:"ORDERS_SERVICE_ADDR" env-default:"localhost:9002"`
	Kitchen   string `yaml:"kitchen" env:"KITCHEN_SERVICE_ADDR" env-default:"localhost:9003"`
	Logistics string `yaml:"logistics" env:"LOGISTICS_SERVICE_ADDR" env-default:"localhost:9004"`
	Chat      string `yaml:"chat" env:"CHAT_SERVICE_ADDR" env-default:"localhost:8089"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("config.yaml", cfg)
	if err != nil {
		err = cleanenv.ReadEnv(cfg)
		if err != nil {
			return nil, fmt.Errorf("config error: %w", err)
		}
	}

	return cfg, nil
}
