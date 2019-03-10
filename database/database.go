package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB // connection handle for the api

func InitDatabase() {
	var err error
	db, err = sql.Open("mysql", "root:root@/sc2hub")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

func Close() {
	err := db.Close()
	if err != nil {
		log.Println(err.Error())
	}
}
