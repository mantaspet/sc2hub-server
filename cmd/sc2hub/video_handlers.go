package main

import (
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strconv"
	"strings"
)

func getPaginatedVideosResponse(videos []*models.Video, cursor int) models.PaginatedVideos {
	var res models.PaginatedVideos
	itemCount := len(videos)
	if itemCount < models.ArticlePageLength+1 {
		res = models.PaginatedVideos{
			Cursor: 0,
			Items:  videos,
		}
	} else {
		res = models.PaginatedVideos{
			Cursor: cursor,
			Items:  videos[:itemCount-1],
		}
	}
	return res
}

func (app *application) getAllVideos(w http.ResponseWriter, r *http.Request) {
	var videos []*models.Video
	var err error

	from := parsePaginationParam(r.URL.Query().Get("from"))

	if r.URL.Query().Get("recent") != "" {
		videos, err = app.videos.SelectRecent()
	} else {
		videos, err = app.videos.SelectPage(models.VideoPageLength, from, r.URL.Query().Get("query"))
	}

	if err != nil {
		app.serverError(w, err)
		return
	}

	res := getPaginatedVideosResponse(videos, from+models.VideoPageLength)
	app.json(w, res)
}

func (app *application) getVideosByPlayer(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(w, r)
	if err != nil {
		return
	}

	from := parsePaginationParam(r.URL.Query().Get("from"))

	videos, err := app.videos.SelectByPlayer(models.VideoPageLength, from, r.URL.Query().Get("query"), id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	res := getPaginatedVideosResponse(videos, from+models.VideoPageLength)
	app.json(w, res)
}

func (app *application) getVideosByCategory(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(w, r)
	if err != nil {
		return
	}

	from := parsePaginationParam(r.URL.Query().Get("from"))

	videos, err := app.videos.SelectByCategory(models.VideoPageLength, from, r.URL.Query().Get("query"), id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	res := getPaginatedVideosResponse(videos, from+models.VideoPageLength)
	app.json(w, res)
}

func (app *application) getEventBroadcasts(w http.ResponseWriter, r *http.Request) {
	id, err := app.parseIDParam(w, r)
	if err != nil {
		return
	}

	videos, err := app.videos.SelectEventBroadcasts(id, r.URL.Query().Get("date"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	res := getPaginatedVideosResponse(videos, 0)
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
