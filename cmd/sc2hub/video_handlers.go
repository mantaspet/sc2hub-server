package main

import (
	"github.com/go-chi/chi"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) getAllVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := app.videos.SelectPage(r.URL.Query().Get("from"), r.URL.Query().Get("query"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	var res models.PaginatedVideos
	itemCount := len(videos)
	if itemCount < models.VideoPageLength {
		res = models.PaginatedVideos{
			Cursor: nil,
			Items:  videos,
		}
	} else {
		res = models.PaginatedVideos{
			Cursor: &videos[itemCount-1].CreatedAt,
			Items:  videos[:itemCount-1],
		}
	}

	app.json(w, res)
}

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

func (app *application) getEventBroadcasts(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	var date string
	if r.URL.Query()["date"] != nil {
		date = r.URL.Query()["date"][0]
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, err)
		return
	}

	res, err := app.videos.SelectEventBroadcasts(id, date)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}

func (app *application) getVideosByPlayer(w http.ResponseWriter, r *http.Request) {
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

	res, err := app.videos.SelectByPlayer(id, query)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}

func (app *application) queryVideoAPIs(w http.ResponseWriter, r *http.Request) {
	channels, err := app.channels.SelectFromAllCategories(0)
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
			videosToInsert = append(videosToInsert, videos...)
		}
	}

	rowCnt, err := app.videos.InsertOrUpdateMany(videosToInsert)
	if err != nil {
		app.serverError(w, err)
		return
	}

	players, err := app.players.SelectAllPlayerIDs()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var playerVideos []models.PlayerVideo
	for _, v := range videosToInsert {
		for _, p := range players {
			if strings.Contains(v.Title, p.PlayerID) {
				playerVideo := models.PlayerVideo{
					PlayerID: p.ID,
					VideoID:  v.ID,
				}
				playerVideos = append(playerVideos, playerVideo)
				break
			}
		}
	}
	_, err = app.players.InsertPlayerVideos(playerVideos)
	if err != nil {
		app.serverError(w, err)
	}

	rowCntStr := strconv.Itoa(int(rowCnt))
	res := "Rows affected: " + rowCntStr

	app.json(w, res)
}
