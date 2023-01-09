package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	// for RunMigrations function
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// for RunMigrations function
	_ "github.com/golang-migrate/migrate/v4/source/file"
)


const migrationsFolder = "file://../../db/migrations"

// RunMigrations performs megrations in specified direction (up or down)
func RunMigrations(dsn string, direction string) error {
	migrate, err := migrate.New(migrationsFolder, dsn)
	if err != nil {
		return fmt.Errorf("cannot init migrate object: %s", err.Error())
	}
	defer func() {
		sourceErr, dbErr := migrate.Close()
	if sourceErr != nil || dbErr != nil {
		log.Fatalf("Cant close migration: sourceErr: %v, dbErr: %v", sourceErr, dbErr)
	}
  	}()
	switch direction {
	case "up":
		return migrate.Up()
	case "down":
		return migrate.Down()
	default:
		return fmt.Errorf("bad direction parameter: %s (only up/down are supported)", direction)
	}
}