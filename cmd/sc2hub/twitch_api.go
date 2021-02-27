package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strings"
	"time"
)

type TwitchVideo struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	CreatedAt    string `json:"created_at"`
	PublishedAt  string `json:"published_at"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Viewable     string `json:"viewable"`
	ViewCount    int    `json:"view_count"`
	Language     string `json:"language"`
	Type         string `json:"type"`
	Duration     string `json:"duration"`
}

type TwitchUser struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	Title           string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
}

type TwitchStream struct {
	UserName     string `json:"user_name"`
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnail_url"`
	ViewerCount  int    `json:"viewer_count"`
}

type TwitchStreamResponse struct {
	Data       []TwitchStream
	Pagination struct {
		Cursor string
	}
}

type StreamResponse struct {
	Items  []TwitchStream
	Cursor string
}

func (app *application) getTwitchAccessToken() error {
	authURL := "https://id.twitch.tv/oauth2/token?grant_type=client_credentials&client_id=" + app.twitchClientId + "&client_secret=" + app.twitchClientSecret
	res, err := app.httpClient.Post(authURL, "application/json", nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	appCredentials := map[string]interface{}{}
	err = json.NewDecoder(res.Body).Decode(&appCredentials)
	if err != nil {
		app.errorLog.Println(err.Error())
		return err
	}

	app.twitchAccessToken = fmt.Sprintf("%v", appCredentials["access_token"])

	return nil
}

func parseTwitchVideoDuration(durationString string) int {
	// twitch duration looks like this: 1h25m30s
	var durationInSeconds int
	duration, err := time.ParseDuration(strings.ToLower(durationString))

	if err != nil {
		durationInSeconds = 0
	} else {
		durationInSeconds = int(duration.Seconds())
	}
	return durationInSeconds
}

// Gets the latest Twitch video data from a given channel
func (app *application) getTwitchVideos(channel *models.Channel) ([]*models.Video, error) {
	url := "https://api.twitch.tv/helix/videos?user_id=" + channel.ID

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+app.twitchAccessToken)
	req.Header.Set("Client-ID", app.twitchClientId)
	res, err := app.httpClient.Do(req)
	if err != nil {
		app.errorLog.Println(err.Error())
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		res, err = app.reauthenticateAndRepeatTwitchRequest(req)
		if err != nil {
			return nil, err
		}
	}

	type Response struct {
		Data       []TwitchVideo
		Pagination interface{}
	}

	var data Response
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	var videos []*models.Video
	for _, v := range data.Data {
		if !app.matchesPattern([]string{v.Title}, channel.IncludePatterns, channel.ExcludePatterns) {
			continue
		}

		createdAt, err := time.Parse("2006-01-02T15:04:05Z", v.PublishedAt)
		if err != nil {
			createdAt = time.Now()
		}

		video := &models.Video{
			ID:              v.ID,
			PlatformID:      1,
			EventCategoryID: channel.EventCategoryID,
			ChannelID:       channel.ID,
			Title:           v.Title,
			Duration:        parseTwitchVideoDuration(v.Duration),
			ThumbnailURL:    v.ThumbnailURL,
			ViewCount:       v.ViewCount,
			Type:            v.Type,
			CreatedAt:       createdAt,
			UpdatedAt:       time.Now(),
		}
		videos = append(videos, video)
	}

	return videos, nil
}

// Gets updated video data of videos already saved in database.
// Returns only a subset of video data, for use with app.videos.UpdateMetadata.
func (app *application) getExistingTwitchVideoData(videos []*models.Video) ([]*models.Video, error) {
	url := "https://api.twitch.tv/helix/videos?"
	for _, v := range videos {
		url += "id=" + v.ID + "&"
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+app.twitchAccessToken)
	res, err := app.httpClient.Do(req)
	if err != nil {
		app.errorLog.Println(err.Error())
		return nil, err
	}

	// Some extra checks here, because based on the return value of this method videos will be deleted.
	// Want to make sure videos that are still present in Twitch do not get deleted.
	if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNotFound) {
		return nil, errors.New("twitch API error")
	}

	if res.StatusCode == http.StatusUnauthorized {
		res, err = app.reauthenticateAndRepeatTwitchRequest(req)
		if err != nil {
			return nil, err
		}
	}

	type Response struct {
		Data       []TwitchVideo
		Pagination interface{}
	}

	var data Response
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	var updatedVideos []*models.Video
	for _, v := range data.Data {
		video := &models.Video{
			ID:           v.ID,
			Title:        v.Title,
			Duration:     parseTwitchVideoDuration(v.Duration),
			ThumbnailURL: v.ThumbnailURL,
			ViewCount:    v.ViewCount,
		}
		updatedVideos = append(updatedVideos, video)
	}
	return updatedVideos, nil
}

func (app *application) getTwitchLiveStreams(query string) (*StreamResponse, error) {
	url := "https://api.twitch.tv/helix/streams?"
	req, _ := http.NewRequest("GET", url+query, nil)
	req.Header.Set("Client-ID", "hmw2ygtkoc9si4001jxq2xmrmc8g99")
	req.Header.Set("Authorization", "Bearer "+app.twitchAccessToken)

	twitchRes, err := app.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if twitchRes.StatusCode == http.StatusUnauthorized {
		twitchRes, err = app.reauthenticateAndRepeatTwitchRequest(req)
		if err != nil {
			return nil, err
		}
	}

	var twitchResDecoded TwitchStreamResponse
	err = json.NewDecoder(twitchRes.Body).Decode(&twitchResDecoded)

	res := StreamResponse{
		Items:  twitchResDecoded.Data,
		Cursor: twitchResDecoded.Pagination.Cursor,
	}

	return &res, err
}

func (app *application) getTwitchChannelDataByLogin(login string, httpClient *http.Client) (models.Channel, error) {
	var channel models.Channel
	url := "https://api.twitch.tv/helix/users?login=" + login
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Client-ID", "hmw2ygtkoc9si4001jxq2xmrmc8g99")
	req.Header.Set("Authorization", "Bearer "+app.twitchAccessToken)
	res, err := httpClient.Do(req)
	if err != nil {
		return channel, err
	}

	type ResponseBody struct {
		Data []TwitchUser
	}
	var resBody ResponseBody
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if len(resBody.Data) == 0 {
		return channel, errors.New("channel does not exist")
	}

	channel = models.Channel{
		ID:              resBody.Data[0].ID,
		PlatformID:      1,
		Login:           resBody.Data[0].Login,
		Title:           resBody.Data[0].Title,
		ProfileImageURL: resBody.Data[0].ProfileImageURL,
	}

	return channel, nil
}

func (app *application) reauthenticateAndRepeatTwitchRequest(req *http.Request) (*http.Response, error) {
	err := app.getTwitchAccessToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+app.twitchAccessToken)
	res, err := app.httpClient.Do(req)
	if err != nil {
		app.errorLog.Println(err.Error())
		return nil, err
	}
	return res, nil
}
