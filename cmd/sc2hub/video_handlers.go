package main

import (
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strconv"
	"time"
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

func (app *application) getVideosFromTwitch(w http.ResponseWriter, r *http.Request) {
	tcs, err := app.twitchChannels.SelectAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var videosToInsert []*models.Video
	for _, tc := range tcs {
		videos, err := app.getTwitchVideos(tc)
		if err != nil {
			app.serverError(w, err)
			return
		}

		for _, v := range videos {
			createdAt, err := time.Parse("2006-01-02T15:04:05Z", v.CreatedAt)
			if err != nil {
				createdAt = time.Now()
			}
			videoID, err := strconv.Atoi(v.ID)
			if err != nil {
				continue
			}
			videoToInsert := &models.Video{
				EventCategoryID: tc.EventCategoryID,
				TwitchID:        videoID,
				ChannelID:       tc.ID,
				Title:           v.Title,
				Duration:        v.Duration,
				CreatedAt:       createdAt,
			}
			videosToInsert = append(videosToInsert, videoToInsert)
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
