package api

import "database/sql"

var DB *sql.DB // connection handle for the api

func InitDatabase() {
	var err error
	DB, err = sql.Open("mysql", "root:root@/sc2hub")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	// Open doesn't open a connection. Validate DSN data:
	err = DB.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}
