package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) getChannelsByCategory(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	res, err := app.twitchChannels.SelectByCategory(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}

func (app *application) addChannelToCategory(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	type Request struct {
		URL string
	}
	var req Request

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	index := strings.Index(req.URL, "twitch.tv/")

	if index > -1 {
		req.URL = req.URL[index+10:]
		index = strings.Index(req.URL, "/")
		if index > -1 {
			req.URL = req.URL[:index]
		}
	} else {
		app.validationError(w, map[string]string{"url": "Must be a valid twitch.tv channel URL"})
		return
	}

	tc, err := app.getChannelDataByLogin(req.URL)
	if err != nil {
		if index := strings.Index(err.Error(), "channel does not exist"); index > -1 {
			app.validationError(w, map[string]string{"url": "Channel does not exist"})
		} else {
			app.serverError(w, err)
		}
		return
	}

	tc.EventCategoryID = id
	res, err := app.twitchChannels.Insert(tc)
	if err != nil {
		if index := strings.Index(err.Error(), "Duplicate entry"); index > -1 {
			app.validationError(w, map[string]string{"url": "This channel is already in database"})
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.json(w, res)
}
