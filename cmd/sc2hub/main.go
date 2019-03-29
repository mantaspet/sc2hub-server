package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"github.com/mantaspet/sc2hub-server/pkg/models/mysql"
	"golang.org/x/crypto/acme/autocert"
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
		InsertMany(events []models.Event) (int64, error)
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
	players interface {
		SelectAllPlayers() ([]*models.Player, error)
		InsertMany(players []models.Player) (int64, error)
	}
}

var (
	flgProduction = false
	flgAddr       = ":9000"
	flgDsn        = ""
)

func parseFlags() {
	flag.BoolVar(&flgProduction, "prod", false, "if true, we start HTTPS server")
	flag.StringVar(&flgAddr, "addr", ":9000", "HTTPS network address")
	flag.StringVar(&flgDsn, "dsn", "root:root@/sc2hub", "MySQL data source name")
	flag.Parse()
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

func main() {
	parseFlags()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	var err error
	db, err := openDB(flgDsn + "?parseTime=true")
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
		players:         &mysql.PlayerModel{DB: db},
	}

	srv := &http.Server{
		Addr:     flgAddr,
		ErrorLog: errorLog,
		Handler:  app.router(),
	}

	if flgProduction {
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("certs"),
		}

		srv.TLSConfig = &tls.Config{GetCertificate: certManager.GetCertificate}
		infoLog.Printf("Starting server on %s", flgAddr)
		go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		err = srv.ListenAndServeTLS("", "")
		if err != nil {
			errorLog.Fatal(err)
		}
	} else {
		infoLog.Printf("Starting server on %s", flgAddr)
		err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
		if err != nil {
			errorLog.Fatal(err)
		}
	}
}
