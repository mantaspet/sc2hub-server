package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strconv"
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

type TwitchChannel struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
}

func (app *application) getTwitchAccessToken() error {
	authURL := "https://id.twitch.tv/oauth2/token?client_secret=7stuc2sc1z5crnrcdtiw9x95cfyqp0&client_id=hmw2ygtkoc9si4001jxq2xmrmc8g99&grant_type=client_credentials"
	res, err := app.httpClient.Post(authURL, "application/json", nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	appCredentials := map[string]interface{}{}
	err = json.NewDecoder(res.Body).Decode(&appCredentials)
	if err != nil {
		return err
	}

	app.twitchAccessToken = fmt.Sprintf("%v", appCredentials["access_token"])

	return nil
}

func (app *application) getTwitchVideos(channel *models.TwitchChannel) ([]TwitchVideo, error) {
	url := "https://api.twitch.tv/helix/videos?user_id=" + strconv.Itoa(channel.TwitchUserID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+app.twitchAccessToken)
	res, err := app.httpClient.Do(req)
	if err != nil {
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

	return data.Data, nil
}

func (app *application) getChannelDataByLogin(login string) (models.TwitchChannel, error) {
	var tc models.TwitchChannel
	url := "https://api.twitch.tv/helix/users?login=" + login
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Client-ID", "hmw2ygtkoc9si4001jxq2xmrmc8g99")
	res, err := app.httpClient.Do(req)
	if err != nil {
		return tc, err
	}

	type Response struct {
		Data []TwitchChannel
	}
	var data Response
	err = json.NewDecoder(res.Body).Decode(&data)
	if len(data.Data) == 0 {
		return tc, errors.New("channel does not exist")
	}

	id, _ := strconv.Atoi(data.Data[0].ID)
	tc = models.TwitchChannel{
		Login:           data.Data[0].Login,
		DisplayName:     data.Data[0].DisplayName,
		TwitchUserID:    id,
		ProfileImageURL: data.Data[0].ProfileImageURL,
	}

	return tc, nil
}

func (app *application) reauthenticateAndRepeatTwitchRequest(req *http.Request) (*http.Response, error) {
	err := app.getTwitchAccessToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+app.twitchAccessToken)
	res, err := app.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
