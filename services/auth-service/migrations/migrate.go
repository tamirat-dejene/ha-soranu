package migrations

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	goose "github.com/pressly/goose/v3"
	authservice "github.com/tamirat-dejene/ha-soranu/services/auth-service"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"go.uber.org/zap"
)

type migrator struct {
	env authservice.Env
}

func NewMigrator(env authservice.Env) *migrator {
	return &migrator{env: env}
}

func (f *migrator) Migrate(ctx context.Context, dir string) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		f.env.DBUser, f.env.DBPassword, f.env.DBHost, f.env.DBPort, f.env.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}
	defer db.Close()

	if err := goose.Up(db, dir); err != nil {
		logger.Error("Failed to run migrations", zap.Error(err))
	}
	return nil
}
