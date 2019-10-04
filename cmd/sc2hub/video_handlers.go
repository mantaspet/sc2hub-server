package main

import (
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	go app.updateVideoMetadata(videos)
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
	go app.updateVideoMetadata(videos)
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
	go app.updateVideoMetadata(videos)
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
	go app.updateVideoMetadata(videos)
	app.json(w, res)
}

func (app *application) initVideoQuerying(w http.ResponseWriter, r *http.Request) {
	res, err := app.queryVideoAPIs()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.json(w, res)
}

func (app *application) queryVideoAPIs() (string, error) {
	channels, err := app.channels.SelectFromAllCategories(0)
	if err != nil {
		return "", err
	}

	var videos []*models.Video
	var videosToInsert []*models.Video
	for _, channel := range channels {
		if channel.PlatformID == 1 {
			videos, err = app.getTwitchVideosByChannel(channel)
		} else if channel.PlatformID == 2 {
			videos, err = app.getYoutubeVideos(channel)
		}

		if err != nil {
			return "", nil
		}

		if len(videos) > 0 {
			videosToInsert = append(videosToInsert, videos...)
		}
	}

	rowCnt, err := app.videos.InsertOrUpdateMany(videosToInsert)
	if err != nil {
		return "", nil
	}

	players, err := app.players.SelectAllPlayerIDsAndIDs()
	if err != nil {
		return "", nil
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
		return "", nil
	}

	rowCntStr := strconv.Itoa(int(rowCnt))
	res := "Rows affected: " + rowCntStr

	return res, nil
}

/**
Get updated metadata (mainly for view count) from video APIs
if it's been 2 hours since the last update for every video.
Current implementation calls this function with "go", to avoid making the user wait.
Might want to rethink this, because it does end up returning outdated info.
*/
func (app *application) updateVideoMetadata(videos []*models.Video) {
	var twitchVideosToUpdate []*models.Video
	for _, v := range videos {
		if time.Now().After(v.UpdatedAt.Add(time.Hour*2)) && v.PlatformID == 1 {
			twitchVideosToUpdate = append(twitchVideosToUpdate, v)
		}
	}

	if len(twitchVideosToUpdate) == 0 {
		return
	}

	updatedVideos, err := app.getTwitchVideos(twitchVideosToUpdate)
	if err != nil {
		app.errorLog.Println("failed to update video metadata: " + err.Error())
	}

	_, err = app.videos.InsertOrUpdateMany(updatedVideos)
	if err != nil {
		app.errorLog.Println("failed to update video metadata: " + err.Error())
	}
}
