package migrations

import (
	"context"
	"database/sql"
	"embed"

	svc "github.com/tamirat-dejene/ha-soranu/services/notification-service"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

//go:embed *.sql
var embedMigrations embed.FS

type Migrator struct {
	env svc.Env
}

func NewMigrator(env svc.Env) *Migrator {
	return &Migrator{env: env}
}

func (m *Migrator) Migrate(ctx context.Context, migrationPath string) error {
	dsn := "postgres://" + m.env.DBUser + ":" + m.env.DBPassword + "@" + m.env.DBHost + ":" + m.env.DBPort + "/" + m.env.DBName + "?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("failed to open database connection", zap.Error(err))
		return err
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("failed to set dialect", zap.Error(err))
		return err
	}

	if err := goose.Up(db, "."); err != nil {
		logger.Error("failed to run migrations", zap.Error(err))
		return err
	}

	logger.Info("migrations applied successfully")
	return nil
}
