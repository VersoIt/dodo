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
	Auth     AuthConfig     `yaml:"auth"`
}

type AppConfig struct {
	Name string `yaml:"name" env:"APP_NAME" env-default:"chat-service"`
	Env  string `yaml:"env" env:"APP_ENV" env-default:"development"`
}

type GRPCConfig struct {
	Port string `yaml:"port" env:"GRPC_PORT" env-default:"9009"`
}

type DatabaseConfig struct {
	URL string `yaml:"url" env:"DATABASE_URL" env-required:"true"`
}

type NatsConfig struct {
	URL string `yaml:"url" env:"NATS_URL" env-default:"nats://localhost:4222"`
}

type AuthConfig struct {
	JWTSecret string `yaml:"jwt_secret" env:"JWT_SECRET" env-default:"super-secret-key"`
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
