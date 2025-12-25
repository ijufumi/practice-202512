package config

import (
	"os"

	"github.com/shopspring/decimal"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	FeeRate    decimal.Decimal
	TaxRate    decimal.Decimal
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "practice"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		FeeRate:    getDecimalEnv("FEE_RATE", "0.04"),
		TaxRate:    getDecimalEnv("TAX_RATE", "0.10"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDecimalEnv(key, defaultValue string) decimal.Decimal {
	value := getEnv(key, defaultValue)
	dec, err := decimal.NewFromString(value)
	if err != nil {
		dec, _ = decimal.NewFromString(defaultValue)
	}
	return dec
}
