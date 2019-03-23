package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mantaspet/sc2hub-server/pkg/models/mysql"
	"log"
	"net/http"
	"os"
)

// TODO inject in handler functions like this:
// func (app *application) home(w http.ResponseWriter, r *http.Request) {
// ...
// app.errorLog.Println(err.Error())
type application struct {
	errorLog        *log.Logger
	infoLog         *log.Logger
	events          *mysql.EventModel
	eventCategories *mysql.EventCategoryModel
}

func main() {
	addr := flag.String("addr", ":9000", "HTTP network address")
	dsn := flag.String("dsn", "root:root@/sc2hub", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	var err error
	db, err := initDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog:        errorLog,
		infoLog:         infoLog,
		events:          &mysql.EventModel{DB: db},
		eventCategories: &mysql.EventCategoryModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.initRouter(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func initDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
