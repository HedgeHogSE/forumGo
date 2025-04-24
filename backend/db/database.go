package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

func SetupDB() {
	var err error
	connStr := "user=postgres password=Sashaezhak2006 dbname=forumDB sslmode=disable"
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	if err = Db.Ping(); err != nil {
		log.Fatal("БД недоступна:", err)
	}
}
