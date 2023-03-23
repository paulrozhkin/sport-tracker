package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/paulrozhkin/sport-tracker/config"
	"go.uber.org/zap"
)

type Store struct {
	Db *sql.DB
}

func NewPostgresStore(cfg *config.Configuration,
	logger *zap.SugaredLogger) (*Store, error) {
	// TODO: заменить на pgxpool
	database := cfg.Database
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		database.DBUser, database.DBPassword, database.DBConnection, database.DBName, database.DBSslMode)
	db, err := sql.Open("postgres",
		psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///Data/Projects/sport-tracker/data/migrations",
		database.DBName, driver)
	if err != nil {
		return nil, err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}
	logger.Info("Database connected")
	return &Store{Db: db}, nil
}
