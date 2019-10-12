package main

import (
	"encoding/json"
	"errors"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type YoutubeVideo struct {
	Id struct {
		VideoId string
	}
	Snippet struct {
		PublishedAt          string
		Title                string
		LiveBroadcastContent string
		Thumbnails           struct {
			Medium struct {
				Url string
			}
		}
	}
}

func (app *application) getYoutubeVideos(channel *models.Channel) ([]*models.Video, error) {
	url := "https://www.googleapis.com/youtube/v3/search" +
		"?key=" + os.Getenv("YOUTUBE_API_KEY") +
		"&channelId=" + channel.ID +
		"&type=video" +
		"&part=snippet,id" +
		"&fields=items(id,snippet(publishedAt,title,liveBroadcastContent,thumbnails(medium(url))))" +
		"&order=date" +
		"&maxResults=50"

	res, err := app.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Items []YoutubeVideo
	}

	var data Response
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	var videos []*models.Video
	for _, v := range data.Items {
		if v.Snippet.LiveBroadcastContent != "none" {
			continue
		}
		match, _ := regexp.MatchString("^[[[:ascii:]]+$", v.Snippet.Title) // to exclude videos with non-english titles
		if match != true {
			continue
		}
		if strings.Contains(strings.ToLower(v.Snippet.Title), channel.Pattern) != true {
			continue
		}
		createdAt, err := time.Parse("2006-01-02T15:04:05.000Z", v.Snippet.PublishedAt)
		if err != nil {
			createdAt = time.Now()
		}
		video := &models.Video{
			ID:              v.Id.VideoId,
			PlatformID:      2,
			EventCategoryID: channel.EventCategoryID,
			ChannelID:       channel.ID,
			Title:           v.Snippet.Title,
			ThumbnailURL:    v.Snippet.Thumbnails.Medium.Url,
			CreatedAt:       createdAt,
		}
		videos = append(videos, video)
	}

	return videos, nil
}

var getYoutubeChannelData = func(login string, id string, httpClient *http.Client) (models.Channel, error) {
	var yc models.Channel
	url := "https://www.googleapis.com/youtube/v3/channels" +
		"?key=" + os.Getenv("YOUTUBE_API_KEY") +
		"&part=id,snippet" +
		"&fields=items(id,snippet(title,customUrl,thumbnails(default)))"

	if login != "" {
		url += "&forUsername=" + login
	} else if id != "" {
		url += "&id=" + id
	} else {
		return yc, errors.New("need to specify either login or id")
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return yc, err
	}

	type ResponseBody struct {
		Items []struct {
			Id      string
			Snippet struct {
				Title      string
				CustomUrl  string
				Thumbnails struct {
					Default struct {
						Url string
					}
				}
			}
		}
	}

	var resBody ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&resBody)
	if len(resBody.Items) == 0 {
		return yc, errors.New("channel does not exist")
	}

	yc = models.Channel{
		ID:              resBody.Items[0].Id,
		Login:           resBody.Items[0].Snippet.CustomUrl,
		PlatformID:      2,
		Title:           resBody.Items[0].Snippet.Title,
		ProfileImageURL: resBody.Items[0].Snippet.Thumbnails.Default.Url,
	}

	return yc, nil
}
