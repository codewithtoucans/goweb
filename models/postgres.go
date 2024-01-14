package models

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DATABASE"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
}

func Open(config PostgresConfig) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), config.String())
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)
}

func migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate %w", err)
	}
	return nil
}

func MigrateFS(migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)
	port, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	sqlStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		port,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_SSLMODE"),
	)
	defer func() { goose.SetBaseFS(nil) }()
	db, err := goose.OpenDBWithDriver("pgx", sqlStr)
	if err != nil {
		return fmt.Errorf("migrate fs %w", err)
	}
	return migrate(db, dir)
}
