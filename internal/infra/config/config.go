package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Server Server
	App    App
	DB     DB
	Auth   Auth
}

type Auth struct {
	JWTSecret  string
	JWTExpiry  time.Duration
	BcryptCost int
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
	_ = godotenv.Load()

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
		Auth: Auth{
			JWTSecret:  envOrDefault("JWT_SECRET", "dev-insecure-change-me"),
			JWTExpiry:  jwtExpiry(),
			BcryptCost: bcryptCost(),
		},
		DB: DB{
			Host:     envOrDefault("DB_HOST", "localhost"),
			Port:     dbPort(),
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

func jwtExpiry() time.Duration {
	raw := envOrDefault("JWT_EXPIRY", "24h")
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}

func bcryptCost() int {
	raw := envOrDefault("BCRYPT_COST", "12")
	cost, err := strconv.Atoi(raw)
	if err != nil || cost < 4 {
		return 12
	}
	return cost
}

func dbPort() string {
	if port := os.Getenv("DB_PORT"); port != "" {
		return port
	}
	if port := os.Getenv("DB_HOST_PORT"); port != "" {
		return port
	}

	return "15432"
}
