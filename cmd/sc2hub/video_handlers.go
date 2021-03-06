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
	channels, err := app.channels.SelectForCrawling(0)
	if err != nil {
		return err.Error(), err
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
			return err.Error(), err
		}

		if len(videos) > 0 {
			videosToInsert = append(videosToInsert, videos...)
		}
	}

	if len(videosToInsert) == 0 {
		return "No videos found", nil
	}

	rowCnt, err := app.videos.InsertOrUpdateMany(videosToInsert)
	if err != nil {
		return err.Error(), err
	}

	players, err := app.players.SelectAllPlayerIDsAndIDs()
	if err != nil {
		return err.Error(), err
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
	if len(playerVideos) > 0 {
		_, err = app.players.InsertPlayerVideos(playerVideos)
		if err != nil {
			return err.Error(), err
		}
	}

	rowCntStr := strconv.Itoa(int(rowCnt))
	res := "Rows affected: " + rowCntStr

	return res, err
}

/**
Get updated metadata (mainly for view count) from video APIs
if it's been 2 hours since the last update for every video.
Current implementation calls this function with "go", to avoid making the user wait.
Might want to rethink this, because it does end up returning outdated info.
*/
func (app *application) updateVideoMetadata(videos []*models.Video) {
	var twitchVideosToUpdate []*models.Video
	var youtubeVideosToUpdate []*models.Video
	for _, v := range videos {
		if time.Now().After(v.UpdatedAt.Add(time.Hour * 2)) {
			if v.PlatformID == 1 {
				twitchVideosToUpdate = append(twitchVideosToUpdate, v)
			} else if v.PlatformID == 2 {
				youtubeVideosToUpdate = append(youtubeVideosToUpdate, v)
			}
		}
	}

	var updatedVideos []*models.Video
	var videosToDelete []*models.Video
	if len(twitchVideosToUpdate) > 0 {
		updatedTwitchVideos, err := app.getExistingTwitchVideoData(twitchVideosToUpdate)
		if err != nil {
			app.errorLog.Println(err.Error())
		} else {
			updatedVideos = append(updatedVideos, updatedTwitchVideos...)
			videosToDelete = getVideosToDelete(twitchVideosToUpdate, updatedTwitchVideos)
		}
	}

	if len(youtubeVideosToUpdate) > 0 {
		updatedYoutubeVideos, err := app.getExistingYoutubeVideoData(youtubeVideosToUpdate)
		if err != nil {
			app.errorLog.Println(err.Error())
		} else {
			updatedVideos = append(updatedVideos, updatedYoutubeVideos...)
			videosToDelete = append(videosToDelete, getVideosToDelete(youtubeVideosToUpdate, updatedYoutubeVideos)...)
		}
	}

	if len(updatedVideos) > 0 {
		err := app.videos.UpdateMetadata(updatedVideos)
		if err != nil {
			app.errorLog.Println(err.Error())
		}
	}

	if len(videosToDelete) > 0 {
		err := app.videos.DeleteMany(videosToDelete)
		if err != nil {
			app.errorLog.Println(err.Error())
		}
	}
}

// Videos that were deleted do not return from Twitch and Youtube.
// We find out which ones by comparing the list that was passed to the API with the list that returned.
func getVideosToDelete(videos []*models.Video, updatedVideos []*models.Video) []*models.Video {
	var videosToDelete []*models.Video
	for _, v := range videos {
		videoWasDeleted := true
		for _, uv := range updatedVideos {
			if v.ID == uv.ID {
				videoWasDeleted = false
				break
			}
		}

		if videoWasDeleted {
			videosToDelete = append(videosToDelete, v)
		}
	}
	return videosToDelete
}
