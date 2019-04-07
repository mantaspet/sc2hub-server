package main

import (
	"encoding/json"
	"fmt"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"strconv"
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

func getTwitchAccessToken() (string, error) {
	authURL := "https://id.twitch.tv/oauth2/token?client_secret=7stuc2sc1z5crnrcdtiw9x95cfyqp0&client_id=hmw2ygtkoc9si4001jxq2xmrmc8g99&grant_type=client_credentials"
	myClient := &http.Client{Timeout: 10 * time.Second}
	res, err := myClient.Post(authURL, "application/json", nil)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	appCredentials := map[string]interface{}{}
	err = json.NewDecoder(res.Body).Decode(&appCredentials)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", appCredentials["access_token"]), nil
}

func getTwitchVideos(channel *models.TwitchChannel, token string) ([]TwitchVideo, error) {
	url := "https://api.twitch.tv/helix/videos?user_id=" + strconv.Itoa(channel.TwitchUserID)
	myClient := &http.Client{Timeout: 10 * time.Second}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := myClient.Do(req)
	if err != nil {
		return nil, err
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
