package main

import (
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strconv"
)

func (app *application) getVideosByCategory(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	var query string
	if r.URL.Query()["query"] != nil {
		query = r.URL.Query()["query"][0]
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	res, err := app.videos.SelectByCategory(id, query)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}

func (app *application) queryVideoAPIs(w http.ResponseWriter, r *http.Request) {
	channels, err := app.channels.SelectFromAllCategories()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var videos []*models.Video
	var videosToInsert []*models.Video
	for _, channel := range channels {
		if channel.PlatformID == 1 {
			videos, err = app.getTwitchVideos(channel)
		} else if channel.PlatformID == 2 {
			videos, err = app.getYoutubeVideos(channel)
		}

		if err != nil {
			app.serverError(w, err)
			return
		}

		if len(videos) > 0 {
			app.videos.InsertOrUpdateMany(videosToInsert)
			videosToInsert = append(videosToInsert, videos...)
		}
	}

	rowCnt, err := app.videos.InsertOrUpdateMany(videosToInsert)
	if err != nil {
		app.serverError(w, err)
		return
	}
	rowCntStr := strconv.Itoa(int(rowCnt))
	res := "Rows affected: " + rowCntStr

	app.json(w, res)
}
