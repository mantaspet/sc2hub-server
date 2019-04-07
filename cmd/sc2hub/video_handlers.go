package main

import (
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

func (app *application) getVideosByCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	res, err := app.videos.SelectByCategory(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}

func (app *application) getVideosFromTwitch(w http.ResponseWriter, r *http.Request) {
	tcs, err := app.twitchChannels.SelectAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	twitchToken, err := getTwitchAccessToken()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var videos []TwitchVideo
	for _, tc := range tcs {
		vids, err := getTwitchVideos(tc, twitchToken)
		if err != nil {
			app.serverError(w, err)
			return
		}
		videos = append(videos, vids...)
	}

	app.json(w, videos)
}
