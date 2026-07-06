package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	migrate "github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/Acauhi99/med-vault/internal/shared/config"
)

func runMigrations(ctx context.Context, cfg config.Config, logger *slog.Logger, args []string) error {
	if len(args) == 0 {
		args = []string{"up"}
	}

	connConfig, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("parse database url: %w", err)
	}

	db := stdlib.OpenDB(*connConfig)
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	driver, err := pgxmigrate.WithInstance(db, &pgxmigrate.Config{})
	if err != nil {
		return fmt.Errorf("init migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file:///app/migrations", "pgx", driver)
	if err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	switch args[0] {
	case "up":
		err = m.Up()
	case "down":
		if len(args) > 1 {
			steps, parseErr := strconv.Atoi(args[1])
			if parseErr != nil || steps <= 0 {
				return fmt.Errorf("invalid down steps: %q", args[1])
			}
			err = m.Steps(-steps)
		} else {
			err = m.Down()
		}
	case "version":
		version, dirty, versionErr := m.Version()
		if versionErr != nil {
			if errors.Is(versionErr, migrate.ErrNilVersion) {
				logger.Info("migration version", "version", "none", "dirty", false)
				return nil
			}
			return fmt.Errorf("read migration version: %w", versionErr)
		}
		logger.Info("migration version", "version", version, "dirty", dirty)
		return nil
	default:
		return fmt.Errorf("unknown migration command %q", args[0])
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	logger.Info("migrations complete", "command", args[0])
	return nil
}
