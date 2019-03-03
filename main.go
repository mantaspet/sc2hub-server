package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mantaspet/sc2hub-server/crawler"
	"log"
	"net/http"
	"time"
)

func main() {
	db := InitDatabase()
	defer db.Close()
	PrintEvents(db)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "sc2hub api\nAvailable endpoints:\nGET /events")
		if err != nil {
			log.Fatal("Unable to return response ", err)
		}
	})

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		events := crawler.TeamliquidEvents()
		eventsJson, err := json.Marshal(events)
		if err != nil {
			log.Fatal("Cannot encode to JSON ", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(eventsJson)
		if err != nil {
			log.Fatal("Unable to return response ", err)
		}
	})

	log.Fatal(http.ListenAndServe(":9000", nil))
}

func InitDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:root@/sc2hub")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return db
}

func PrintEvents(db *sql.DB) {
	var (
		id       int
		title    string
		stage    string
		startsAt string
	)
	start := time.Now()
	defer log.Printf("Query successful. Elapsed time: %v\n", time.Since(start))
	rows, err := db.Query("SELECT id, title, stage, starts_at FROM events")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &title, &stage, &startsAt)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("id: ", id, "title: ", title, ", stage: ", stage, ", starts at: ", startsAt)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
