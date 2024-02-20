package main

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"time"
)

const (
	_defaultAttempts = 10
	_defaultTimeout  = time.Second
)

func main() {
	databaseConnectionUri := os.Getenv("PG_URL")
	if len(databaseConnectionUri) == 0 {
		log.Fatalf("migrate: environment variables not declared")
	}

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseConnectionUri)
		if err == nil {
			break
		}

		log.Printf("migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("migrate: no change")
		return
	}
	if err != nil {
		log.Fatalf("migrate: up error: %s", err)
	}

	log.Printf("migrate: up success")
}
