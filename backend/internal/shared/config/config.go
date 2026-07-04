package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env             string        `envconfig:"APP_ENV" default:"development"`
	HTTPAddr        string        `envconfig:"HTTP_ADDR" default:":8080"`
	DatabaseURL     string        `envconfig:"DATABASE_URL"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
	ReadTimeout     time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"15s"`
	WriteTimeout    time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"15s"`
	IdleTimeout     time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" default:"60s"`
	RequestIDHeader string        `envconfig:"REQUEST_ID_HEADER" default:"X-Request-Id"`
	JWTSecret       string        `envconfig:"JWT_SECRET" default:"dev-secret-change-in-production"`
	JWTAccessTTL    time.Duration `envconfig:"JWT_ACCESS_TTL" default:"15m"`
	JWTRefreshTTL   time.Duration `envconfig:"JWT_REFRESH_TTL" default:"168h"`
	JWTTempTTL      time.Duration `envconfig:"JWT_TEMP_TTL" default:"5m"`
	BcryptCost      int           `envconfig:"BCRYPT_COST" default:"12"`
	S3Bucket        string        `envconfig:"S3_BUCKET" default:"med-vault-dev"`
	AWSRegion       string        `envconfig:"AWS_REGION" default:"us-east-1"`
}

func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}

	if cfg.Env == "production" && cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required in production")
	}

	return cfg, nil
}
