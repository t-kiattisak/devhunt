package infrastructure

import (
	"context"
	config "devhunt/configs"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewPostgresDB() *pgxpool.Pool {
	_ = godotenv.Load()
	dsn := config.Get().DBURL

	err := runMigrations(dsn)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse DB config: %v", err)
	}

	dbpool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}

	fmt.Println("Connected to PostgreSQL via pgxpool")
	return dbpool
}

func runMigrations(dsn string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	migrationsPath := filepath.Join(wd, "internal/infrastructure/migrations")
	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}
