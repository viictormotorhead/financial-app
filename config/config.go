package config

import "os"

type Configuration struct {
	Server Server
	App    App
	DB     DB
}

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type Server struct {
	Port     string
	BasePath string
}

type App struct {
	ServiceName string
}

var cfg Configuration

func Load() error {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = ""
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "financial-app"
	}

	cfg = Configuration{
		Server: Server{
			Port:     port,
			BasePath: basePath,
		},
		App: App{ServiceName: serviceName},
		DB: DB{
			Host:     envOrDefault("DB_HOST", "localhost"),
			Port:     envOrDefault("DB_PORT", "5432"),
			User:     envOrDefault("DB_USER", "financial"),
			Password: envOrDefault("DB_PASSWORD", "financial"),
			Name:     envOrDefault("DB_NAME", "financial_app"),
			SSLMode:  envOrDefault("DB_SSLMODE", "disable"),
		},
	}

	return nil
}

func Config() Configuration {
	return cfg
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
