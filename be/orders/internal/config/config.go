package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	Database DatabaseConfig `yaml:"database"`
	Services ServicesConfig `yaml:"services"`
}

type AppConfig struct {
	Name string `yaml:"name" env:"APP_NAME" env-default:"orders-service"`
	Env  string `yaml:"env" env:"APP_ENV" env-default:"development"`
}

type GRPCConfig struct {
	Port string `yaml:"port" env:"GRPC_PORT" env-default:"9002"`
}

type DatabaseConfig struct {
	URL string `yaml:"url" env:"DATABASE_URL" env-required:"true"`
}

type ServicesConfig struct {
	Catalog      string `yaml:"catalog" env:"CATALOG_SERVICE_ADDR" env-default:"localhost:9001"`
	Kitchen      string `yaml:"kitchen" env:"KITCHEN_SERVICE_ADDR" env-default:"localhost:9003"`
	Logistics    string `yaml:"logistics" env:"LOGISTICS_SERVICE_ADDR" env-default:"localhost:9004"`
	Treasury     string `yaml:"treasury" env:"TREASURY_SERVICE_ADDR" env-default:"localhost:9006"`
	Notification string `yaml:"notification" env:"NOTIFICATION_SERVICE_ADDR" env-default:"localhost:9005"`
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
