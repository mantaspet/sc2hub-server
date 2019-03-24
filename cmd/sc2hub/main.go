package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"github.com/mantaspet/sc2hub-server/pkg/models/mysql"
	"log"
	"net/http"
	"os"
)

type application struct {
	db       *sql.DB // TODO find a better solution. This is used only in pkg validators SQLUnique function
	errorLog *log.Logger
	infoLog  *log.Logger
	events   interface {
		SelectInDateRange(dateFrom string, dateTo string) ([]*models.Event, error)
		Insert(events []*models.Event) (int64, error)
	}
	eventCategories interface {
		SelectAll() ([]*models.EventCategory, error)
		Insert(ec models.EventCategory) (*models.EventCategory, error)
		Update(id string, ec models.EventCategory) (*models.EventCategory, error)
		Delete(id string) error
		UpdatePriorities(id int, newPrio int) error
		AssignToEvents(events []*models.Event) ([]*models.Event, error)
		LoadOnEvents(events []*models.Event) ([]*models.Event, error)
	}
}

func main() {
	addr := flag.String("addr", ":9000", "HTTP network address")
	dsn := flag.String("dsn", "root:root@/sc2hub", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	var err error
	db, err := openDB(*dsn + "?parseTime=true")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		db:              db,
		errorLog:        errorLog,
		infoLog:         infoLog,
		events:          &mysql.EventModel{DB: db},
		eventCategories: &mysql.EventCategoryModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.router(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
