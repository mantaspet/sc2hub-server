package main

import (
	"time"
)

func (app *application) initScheduler() {
	hourlyTicker := time.NewTicker(1 * time.Hour)
	go func() {
		for range hourlyTicker.C {
			go app.initQueryVideos()
			go app.initCrawlEvents()
			go app.initCrawlArticles()
		}
	}()

	dailyTicker := time.NewTicker(24 * time.Hour)
	go func() {
		for range dailyTicker.C {
			go app.initCrawlPlayers()
		}
	}()
}

func (app *application) initQueryVideos() {
	res, err := app.queryVideoAPIs()
	if err != nil {
		app.errorLog.Println(err)
	} else {
		app.infoLog.Println("Queried video APIs. " + res)
	}
}

func (app *application) initCrawlEvents() {
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	nextMonth := time.Now().AddDate(0, 1, 0).UTC().Format("01")

	res, err := app.crawlEvents(year, month)
	if err != nil {
		app.errorLog.Println(err)
	} else {
		app.infoLog.Println("Queried " + month + " month events. " + res)
	}

	time.Sleep(1000 * time.Millisecond)
	res, err = app.crawlEvents(year, nextMonth)
	if err != nil {
		app.errorLog.Println(err)
	} else {
		app.infoLog.Println("Queried " + nextMonth + " month events. " + res)
	}
}

func (app *application) initCrawlArticles() {
	res, err := app.crawlArticles()
	if err != nil {
		app.errorLog.Println(err)
	} else {
		app.infoLog.Println("Crawled articles. " + res)
	}
}

func (app *application) initCrawlPlayers() {
	regions := [4]string{"Europe", "US", "Asia", "Korea"}
	for _, r := range regions {
		res, err := app.crawlPlayers(r)
		if err != nil {
			app.errorLog.Println(err)
		} else {
			app.infoLog.Println("Crawled players from + " + r + ". " + res)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
