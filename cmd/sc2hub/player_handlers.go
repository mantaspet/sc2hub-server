package main

import (
	"errors"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/crawlers"
	"net/http"
	"strconv"
	"time"
)

func (app *application) getAllPlayers(w http.ResponseWriter, r *http.Request) {
	events, err := app.players.SelectAllPlayers()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, events)
}

func (app *application) crawlPlayers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer fmt.Printf("Successfully crawled liquipedia players. Elapsed time: %v\n", time.Since(start))
	regionFound := false
	region := r.URL.Query().Get("region")
	regions := [4]string{"Europe", "US", "Asia", "Korea"}
	for _, r := range regions {
		if r == region {
			regionFound = true
			break
		}
	}

	if !regionFound {
		app.clientError(w, http.StatusBadRequest, errors.New("need to specify one of these regions: Europe, US, Asia, Korea"))
		return
	}

	players, err := crawlers.LiquipediaPlayers(region)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if len(players) == 0 {
		app.json(w, "No players found")
		return
	}

	rowCnt, err := app.players.InsertMany(players)
	if err != nil {
		app.serverError(w, err)
		return
	}
	rowCntStr := strconv.Itoa(int(rowCnt))
	res := "Rows affected: " + rowCntStr

	app.json(w, res)
}
