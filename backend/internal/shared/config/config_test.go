package config

import "testing"

func TestLoadComposesDatabaseURL(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("DB_HOST", "db.example.internal")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "medvault")
	t.Setenv("DB_USERNAME", "medvault")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("JWT_SECRET", "jwt-secret")
	t.Setenv("S3_BUCKET", "bucket")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.DatabaseURL == "" {
		t.Fatal("expected database url to be composed")
	}
	if got := cfg.DBSSLMode; got != "require" {
		t.Fatalf("db ssl mode = %q, want require", got)
	}
}
