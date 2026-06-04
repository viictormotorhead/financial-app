package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/viictormotorhead/financial-app/internal/infra/clients/db"
	"github.com/viictormotorhead/financial-app/internal/infra/config"
)

func main() {
	migrationsPath := flag.String("path", "migrations", "directory with SQL migration files")
	flag.Parse()

	if err := config.Load(); err != nil {
		log.Fatalf("load config: %v", err)
	}

	command, args := parseCommand(flag.Args())
	if command == "" {
		printUsage()
		os.Exit(1)
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", *migrationsPath),
		db.MigrationURL(),
	)
	if err != nil {
		log.Fatalf("create migrator: %v", err)
	}
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			log.Printf("close migration source: %v", sourceErr)
		}
		if dbErr != nil {
			log.Printf("close migration database: %v", dbErr)
		}
	}()

	if err := runCommand(m, command, args); err != nil {
		log.Fatal(err)
	}
}

func parseCommand(args []string) (string, []string) {
	if len(args) == 0 {
		return "", nil
	}

	return args[0], args[1:]
}

func runCommand(m *migrate.Migrate, command string, args []string) error {
	switch command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("apply migrations: %w", err)
		}
		log.Println("migrations applied")
		return nil
	case "down":
		steps := 1
		if len(args) > 0 {
			parsed, err := strconv.Atoi(args[0])
			if err != nil || parsed < 1 {
				return fmt.Errorf("invalid down steps: %q", args[0])
			}
			steps = parsed
		}

		if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("rollback migrations: %w", err)
		}
		log.Printf("rolled back %d migration(s)", steps)
		return nil
	case "force":
		if len(args) != 1 {
			return errors.New("force requires a version number")
		}

		version, err := strconv.Atoi(args[0])
		if err != nil || version < 0 {
			return fmt.Errorf("invalid force version: %q", args[0])
		}

		if err := m.Force(version); err != nil {
			return fmt.Errorf("force migration version: %w", err)
		}

		log.Printf("forced migration version to %d", version)
		return nil
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				log.Println("no migrations applied yet")
				return nil
			}
			return fmt.Errorf("read migration version: %w", err)
		}

		log.Printf("current version: %d (dirty=%t)", version, dirty)
		return nil
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}

func printUsage() {
	log.Println(`usage: go run ./cmd/migrate [flags] <command>

commands:
  up                 apply all pending migrations
  down [N]           rollback N migrations (default: 1)
  force VERSION      set schema_migrations version without running SQL
  version            print current migration version`)
}
