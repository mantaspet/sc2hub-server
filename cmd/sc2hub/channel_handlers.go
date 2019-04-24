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

	res, err := app.channels.SelectByCategory(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}

// checks if URL points to a valid twitch or youtube channel,
// fetches it's data from twitch or youtube api,
// stores it inside our database
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

	twitchIndex := strings.Index(req.URL, "twitch.tv/")
	youtubeIndex := strings.Index(req.URL, "youtube.com/user/")

	var channel models.Channel
	if twitchIndex > -1 {
		req.URL = req.URL[twitchIndex+10:]
		twitchIndex = strings.Index(req.URL, "/")
		if twitchIndex > -1 {
			req.URL = req.URL[:twitchIndex]
		}
		channel, err = app.getChannelDataByLogin(req.URL)
	} else if youtubeIndex > -1 {
		req.URL = req.URL[youtubeIndex+17:]
		youtubeIndex = strings.Index(req.URL, "/")
		if youtubeIndex > -1 {
			req.URL = req.URL[:youtubeIndex]
		}
		channel, err = app.getYoutubeChannelDataByLogin(req.URL)
	} else {
		app.validationError(w, map[string]string{"url": "Must be a valid twitch.tv or youtube.com channel URL"})
		return
	}

	if err != nil {
		if index := strings.Index(err.Error(), "channel does not exist"); index > -1 {
			app.validationError(w, map[string]string{"url": "Channel does not exist"})
		} else {
			app.serverError(w, err)
		}
		return
	}

	channel.Login = req.URL
	res, err := app.channels.Insert(channel, id)
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

func (app *application) deleteCategoryChannel(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	if channelID == "" {
		app.clientError(w, http.StatusBadRequest, errors.New("must specify channel ID"))
		return
	}

	categoryID, err := strconv.Atoi(chi.URLParam(r, "categoryID"))
	if err != nil || categoryID < 1 {
		app.clientError(w, http.StatusBadRequest, errors.New("must specify a valid category ID"))
		return
	}

	err = app.channels.DeleteFromCategory(channelID, categoryID)
	if err == models.ErrNotFound {
		app.clientError(w, http.StatusNotFound, errors.New("no channel found in category"))
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
