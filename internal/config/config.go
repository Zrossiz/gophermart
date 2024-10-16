package config

import (
	"errors"
	"flag"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RunAddress            string
	DBDSN                 string
	AcccrualSystemAddress string
	AccessTokenSecret     string
	RefreshTokenSecret    string
	LogLevel              string
	Cost                  int
	AutoMigrate           string
}

var AppConfig *Config

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) (int, error) {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, nil
		}
		return intValue, nil
	}
	return defaultValue, nil
}

func Init() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	flag.StringVar(&cfg.RunAddress, "a", "", "server address")
	flag.StringVar(&cfg.DBDSN, "d", "", "db dsn")
	flag.StringVar(&cfg.AcccrualSystemAddress, "r", "0.0.0.0:1234", "accrual system address")
	flag.StringVar(&cfg.AccessTokenSecret, "as", "access_secret", "access token secret")
	flag.StringVar(&cfg.RefreshTokenSecret, "rs", "refresh_secret", "refresh token secret")
	flag.StringVar(&cfg.LogLevel, "l", "WARN", "log level")
	flag.IntVar(&cfg.Cost, "s", 4, "cost for hash password")
	flag.Parse()

	cfg.AutoMigrate = getEnvOrDefault("AUTO_MIGRATE", "")
	cfg.RunAddress = getEnvOrDefault("RUN_ADDRESS", cfg.RunAddress)
	cfg.DBDSN = getEnvOrDefault("DATABASE_URI", cfg.DBDSN)
	cfg.AcccrualSystemAddress = getEnvOrDefault("ACCRUAL_SYSTEM_ADDRESS", cfg.AcccrualSystemAddress)
	cfg.AccessTokenSecret = getEnvOrDefault("ACCESS_TOKEN_SECRET", cfg.AccessTokenSecret)
	cfg.RefreshTokenSecret = getEnvOrDefault("REFRESH_TOKEN_SECRET", cfg.RefreshTokenSecret)
	cfg.LogLevel = getEnvOrDefault("LOG_LEVEL", cfg.LogLevel)
	cost, err := getIntEnvOrDefault("COST", 4)
	if err != nil {
		return nil, errors.New("invalID cost value")
	}
	cfg.Cost = cost

	if cfg.AccessTokenSecret == "" {
		return nil, errors.New("access token secret not provIDed")
	}
	if cfg.RefreshTokenSecret == "" {
		return nil, errors.New("refresh token secret not provIDed")
	}
	if cfg.AcccrualSystemAddress == "" {
		return nil, errors.New("accrual system address not provIDed")
	}
	if cfg.DBDSN == "" {
		return nil, errors.New("db uri not provIDed")
	}
	if cfg.RunAddress == "" {
		return nil, errors.New("run address not provIDed")
	}

	AppConfig = cfg

	return cfg, nil
}
