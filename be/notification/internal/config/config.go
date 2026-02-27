package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	Database DatabaseConfig `yaml:"database"`
	Nats     NatsConfig     `yaml:"nats"`
}

type AppConfig struct {
	Name string `yaml:"name" env:"APP_NAME" env-default:"notification-service"`
	Env  string `yaml:"env" env:"APP_ENV" env-default:"development"`
}

type GRPCConfig struct {
	Port string `yaml:"port" env:"GRPC_PORT" env-default:"9005"`
}

type DatabaseConfig struct {
	URL string `yaml:"url" env:"DATABASE_URL" env-required:"true"`
}

type NatsConfig struct {
	URL string `yaml:"url" env:"NATS_URL" env-default:"nats://nats:4222"`
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
