package migrations

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	goose "github.com/pressly/goose/v3"
	svc "github.com/tamirat-dejene/ha-soranu/services/restaurant-service"
)

type migrator struct {
	env svc.Env
}

func NewMigrator(env svc.Env) *migrator {
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
		return err
	}
	return nil
}
