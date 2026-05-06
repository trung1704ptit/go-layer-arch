package main

import (
	"log"
	"net/url"
	"os"

	"github.com/go-layer-arch/internal/initializers"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	password := url.QueryEscape(config.DBUserPassword)
	databaseURL := "postgres://" + config.DBUserName + ":" + password + "@" + config.DBHost + ":" + config.DBPort + "/" + config.DBName + "?sslmode=disable"

	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		log.Fatal("🚀 Could not initialize migration runner", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Println("⚠️ Could not close migration source:", srcErr)
		}
		if dbErr != nil {
			log.Println("⚠️ Could not close migration database:", dbErr)
		}
	}()

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		log.Fatalf("unsupported migration command: %s (use up or down)", cmd)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("🚀 Migration failed", err)
	}

	log.Printf("👍 Migration %s complete", cmd)
}
