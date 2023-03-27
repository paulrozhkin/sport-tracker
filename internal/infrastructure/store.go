package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paulrozhkin/sport-tracker/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Store struct {
	Pool *pgxpool.Pool
}

func NewPostgresStore(lc fx.Lifecycle,
	cfg *config.Configuration,
	logger *zap.SugaredLogger) (*Store, error) {
	logger.Info("Start initialize database")
	database := cfg.Database
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		database.DBUser, database.DBPassword, database.DBConnection, database.DBName, database.DBSslMode)
	err := migrateDatabase("file:///Data/Projects/sport-tracker/data/migrations", psqlInfo, logger)
	if err != nil {
		logger.Error("Failed to create migrator due to", zap.Error(err))
		return nil, err
	}
	store := new(Store)
	logger.Info("Database initialized")
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			pool, connectErr := pgxpool.New(ctx, psqlInfo)
			if connectErr != nil {
				logger.Error("Failed to connect to database", zap.Error(connectErr))
				return connectErr
			}
			pingError := pool.Ping(ctx)
			if pingError != nil {
				logger.Error("Failed to ping database", zap.Error(pingError))
				return err
			}
			logger.Info("Database connected")
			store.Pool = pool
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if store.Pool != nil {
				logger.Info("Database closed")
				store.Pool.Close()
			}
			return nil
		},
	})
	return store, nil
}

func migrateDatabase(pathToMigrations, connectionUrl string, logger *zap.SugaredLogger) error {
	m, err := migrate.New(
		pathToMigrations,
		connectionUrl)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
		logger.Info("Database no need migrations")
	} else {
		logger.Info("Database migrated to new version")
	}
	return nil
}
