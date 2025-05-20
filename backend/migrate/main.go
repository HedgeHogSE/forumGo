package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	command = flag.String("command", "", "Команда для выполнения (up, down)")
	service = flag.String("service", "", "Сервис (auth или forum)")
)

func main() {
	flag.Parse()

	if *command == "" {
		os.Exit(1)
	}

	if *service == "" {
		log.Fatal("Необходимо указать сервис (--service)")
	}

	if *service != "auth" && *service != "forum" {
		log.Fatal("Сервис должен быть 'auth' или 'forum'")
	}

	dbURL := fmt.Sprint("postgres://postgres:Sashaezhak2006@localhost:5432/forumDB?sslmode=disable")
	migrationsPath := fmt.Sprintf("file://backend/%s/migrations", *service)

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Ошибка создания миграции: %v", err)
	}
	defer m.Close()

	switch *command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Ошибка применения миграций: %v", err)
		}
		fmt.Printf("Миграции для %s успешно применены\n", *service)

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Ошибка отката миграций: %v", err)
		}
		fmt.Printf("Миграции для %s успешно откачены\n", *service)

	default:
		os.Exit(1)
	}
}
