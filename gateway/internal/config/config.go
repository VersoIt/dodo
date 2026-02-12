package config

import "os"

type Config struct {
	AppPort        string
	LogLevel       string
	AuthService    string
	CatalogService string
	OrdersService  string
		KitchenService   string
		LogisticsService string
		AnalyticsService string
	}
	
	func NewConfig() *Config {
		getEnv := func(key, fallback string) string {
			if value, ok := os.LookupEnv(key); ok {
				return value
			}
			return fallback
		}
	
		return &Config{
			AppPort:          getEnv("APP_PORT", "8080"),
			LogLevel:         getEnv("LOG_LEVEL", "info"),
			AuthService:      getEnv("AUTH_SERVICE_ADDR", "localhost:50051"),
			CatalogService:   getEnv("CATALOG_SERVICE_ADDR", "localhost:50052"),
			OrdersService:    getEnv("ORDERS_SERVICE_ADDR", "localhost:50053"),
			KitchenService:   getEnv("KITCHEN_SERVICE_ADDR", "localhost:50054"),
			LogisticsService: getEnv("LOGISTICS_SERVICE_ADDR", "localhost:50055"),
			AnalyticsService: getEnv("ANALYTICS_SERVICE_ADDR", "localhost:50056"),
		}
	}