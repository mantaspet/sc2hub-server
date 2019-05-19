package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mantaspet/sc2hub-server/pkg/models/mysql"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	flgProduction   = false
	flgAddr         = ":443"
	flgDsn          = ""
	flgOrigin       = ""
	flgClientSecret = ""
)

func parseFlags() {
	flag.BoolVar(&flgProduction, "prod", false, "if true, we start HTTPS server")
	flag.StringVar(&flgAddr, "addr", ":443", "HTTPS network address")
	flag.StringVar(&flgDsn, "dsn", "root:root@/sc2hub", "MySQL data source name")
	flag.StringVar(&flgOrigin, "origin", "http://localhost:4200", "client origin")
	flag.StringVar(&flgClientSecret, "secret", "", "JWT auth client secret")
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
	if flgClientSecret == "" {
		log.Fatal("must specify a value for flag 'secret'")
	}
	signingKey = []byte(flgClientSecret)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	var err error
	db, err := openDB(flgDsn + "?parseTime=true")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	httpClient := &http.Client{Timeout: 10 * time.Second}

	app := &application{
		db:              db,
		httpClient:      httpClient,
		origin:          flgOrigin,
		errorLog:        errorLog,
		infoLog:         infoLog,
		events:          &mysql.EventModel{DB: db},
		eventCategories: &mysql.EventCategoryModel{DB: db},
		players:         &mysql.PlayerModel{DB: db},
		articles:        &mysql.ArticleModel{DB: db},
		videos:          &mysql.VideoModel{DB: db},
		channels:        &mysql.ChannelModel{DB: db},
		users:           &mysql.UserModel{DB: db},
	}

	err = app.getTwitchAccessToken()
	if err != nil {
		app.errorLog.Println(err.Error())
		return
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
		err = srv.ListenAndServe()
		if err != nil {
			errorLog.Fatal(err)
		}
	}
}
