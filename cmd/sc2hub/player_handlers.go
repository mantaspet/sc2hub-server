package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/pkg/crawlers"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strconv"
	"time"
)

func (app *application) getAllPlayers(w http.ResponseWriter, r *http.Request) {
	fromParam := r.URL.Query().Get("from")
	from, err := strconv.Atoi(fromParam)
	if err != nil && fromParam != "" {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	players, err := app.players.SelectPage(from, r.URL.Query().Get("query"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	var res models.PaginatedPlayers
	itemCount := len(players)
	if itemCount < models.PlayerPageLength {
		res = models.PaginatedPlayers{
			Cursor: nil,
			Items:  players,
		}
	} else {
		res = models.PaginatedPlayers{
			Cursor: &players[itemCount-1].ID,
			Items:  players[:itemCount-1],
		}
	}

	app.json(w, res)
}

func (app *application) getPlayer(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	res, err := app.players.SelectOne(id)
	if err == models.ErrNotFound {
		app.clientError(w, http.StatusNotFound, errors.New("player with a specified ID does not exist"))
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
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
