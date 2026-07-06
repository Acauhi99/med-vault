package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env             string        `envconfig:"APP_ENV" default:"development"`
	HTTPAddr        string        `envconfig:"HTTP_ADDR" default:":8080"`
	DatabaseURL     string        `envconfig:"DATABASE_URL"`
	DBHost          string        `envconfig:"DB_HOST"`
	DBPort          int           `envconfig:"DB_PORT" default:"5432"`
	DBName          string        `envconfig:"DB_NAME"`
	DBUsername      string        `envconfig:"DB_USERNAME"`
	DBPassword      string        `envconfig:"DB_PASSWORD"`
	DBSSLMode       string        `envconfig:"DB_SSLMODE" default:"require"`
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

	if cfg.DatabaseURL == "" {
		derivedURL, err := cfg.composeDatabaseURL()
		if err != nil {
			return Config{}, err
		}
		cfg.DatabaseURL = derivedURL
	}

	if cfg.Env == "production" && cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required in production")
	}

	return cfg, nil
}

func (c Config) composeDatabaseURL() (string, error) {
	if c.DBHost == "" || c.DBName == "" || c.DBUsername == "" || c.DBPassword == "" {
		return "", nil
	}

	port := c.DBPort
	if port <= 0 {
		port = 5432
	}
	sslMode := c.DBSSLMode
	if strings.TrimSpace(sslMode) == "" {
		sslMode = "require"
	}

	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.DBUsername, c.DBPassword),
		Host:     netHost(c.DBHost, port),
		Path:     "/" + c.DBName,
		RawQuery: "sslmode=" + url.QueryEscape(sslMode),
	}).String(), nil
}

func netHost(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}
