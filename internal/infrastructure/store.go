package infrastructure

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/paulrozhkin/sport-tracker/config"
)

type Store struct {
	db *sql.DB
}

func CreateAndMigrate(config *config.DatabaseConfigurations) (*Store, error) {
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.DBUser, config.DBPassword, config.DBConnection, config.DBName, config.DBSslMode)
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
		config.DBName, driver)
	if err != nil {
		return nil, err
	}
	err = m.Up()
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}
