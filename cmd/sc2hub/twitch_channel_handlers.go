package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/pkg/models"
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

func (app *application) deleteChannel(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
	}

	err = app.twitchChannels.Delete(id)
	if err == models.ErrNotFound {
		app.clientError(w, http.StatusNotFound, errors.New("twitch channel with a specified ID does not exist"))
		return
	} else if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint fails") {
			app.clientError(w, http.StatusConflict, err)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.json(w, "channel was deleted")
}

func (app *application) getTwitchAppAccessToken(w http.ResponseWriter, r *http.Request) {
	app.json(w, app.twitchAccessToken)
}
