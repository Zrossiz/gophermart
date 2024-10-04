package config

import (
	"errors"
	"flag"
	"os"
)

type Config struct {
	RunAddress            string
	DBDSN                 string
	AcccrualSystemAddress string
	AccessTokenSecret     string
	RefreshTokenSecret    string
}

var AppConfig *Config

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Init() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.RunAddress, "a", "", "server address")
	flag.StringVar(&cfg.DBDSN, "d", "", "db dsn")
	flag.StringVar(&cfg.AcccrualSystemAddress, "r", "", "accrual system address")
	flag.StringVar(&cfg.AccessTokenSecret, "as", "", "access token secret")
	flag.StringVar(&cfg.RefreshTokenSecret, "rs", "", "refresh token secret")
	flag.Parse()

	cfg.RunAddress = getEnvOrDefault("RUN_ADDRESS", cfg.RunAddress)
	cfg.DBDSN = getEnvOrDefault("DATABASE_URI", cfg.DBDSN)
	cfg.AcccrualSystemAddress = getEnvOrDefault("ACCRUAL_SYSTEM_ADDRESS", cfg.AcccrualSystemAddress)
	cfg.AccessTokenSecret = getEnvOrDefault("ACCESS_TOKEN_SECRET", cfg.AccessTokenSecret)
	cfg.RefreshTokenSecret = getEnvOrDefault("REFRESH_TOKEN_SECRET", cfg.RefreshTokenSecret)

	if cfg.AccessTokenSecret == "" {
		return nil, errors.New("access token secret not provided")
	}
	if cfg.RefreshTokenSecret == "" {
		return nil, errors.New("refresh token secret not provided")
	}
	if cfg.AcccrualSystemAddress == "" {
		return nil, errors.New("accrual system address not provided")
	}
	if cfg.DBDSN == "" {
		return nil, errors.New("db uri not provided")
	}
	if cfg.RunAddress == "" {
		return nil, errors.New("run address not provided")
	}

	AppConfig = cfg

	return cfg, nil
}
