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
	flgProduction         = false
	flgAddr               = ":443"
	flgDsn                = ""
	flgAppOrigin          = ""
	flgAdminOrigin        = ""
	flgClientSecret       = ""
	flgYoutubeApiKey      = ""
	flgTwitchClientId     = ""
	flgTwitchClientSecret = ""
)

func parseFlags() {
	flag.BoolVar(&flgProduction, "prod", false, "if true, we start HTTPS server")
	flag.StringVar(&flgAddr, "addr", ":443", "HTTPS network address")
	flag.StringVar(&flgDsn, "dsn", "root:root@/sc2hub", "MySQL data source name")
	flag.StringVar(&flgAppOrigin, "appOrigin", "http://localhost:3000", "app client origin")
	flag.StringVar(&flgAdminOrigin, "adminOrigin", "http://localhost:8080", "admin client origin")
	flag.StringVar(&flgClientSecret, "secret", "", "JWT auth client secret")
	flag.StringVar(&flgYoutubeApiKey, "youtube_key", "", "YouTube data API v3 key")
	flag.StringVar(&flgTwitchClientId, "twitchClientId", "", "Twitch app client id")
	flag.StringVar(&flgTwitchClientSecret, "twitchClientSecret", "", "Twitch app client secret")
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

	db, err := openDB(flgDsn + "?parseTime=true")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	httpClient := &http.Client{Timeout: 10 * time.Second}

	app := &application{
		db:                 db,
		httpClient:         httpClient,
		appOrigin:          flgAppOrigin,
		adminOrigin:        flgAdminOrigin,
		twitchClientId:     flgTwitchClientId,
		twitchClientSecret: flgTwitchClientSecret,
		errorLog:           errorLog,
		infoLog:            infoLog,
		events:             &mysql.EventModel{DB: db},
		eventCategories:    &mysql.EventCategoryModel{DB: db},
		players:            &mysql.PlayerModel{DB: db},
		articles:           &mysql.ArticleModel{DB: db},
		videos:             &mysql.VideoModel{DB: db},
		channels:           &mysql.ChannelModel{DB: db},
		users:              &mysql.UserModel{DB: db},
		twitchGameId:       490422,
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

	app.initScheduler()

	if flgProduction {
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("certs"),
		}

		srv.TLSConfig = &tls.Config{GetCertificate: certManager.GetCertificate}
		infoLog.Printf("Starting server on %s", flgAddr)
		go http.ListenAndServe(":81", certManager.HTTPHandler(nil))
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
