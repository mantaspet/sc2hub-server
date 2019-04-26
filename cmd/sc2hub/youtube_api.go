package main

import (
	"encoding/json"
	"errors"
	"github.com/mantaspet/sc2hub-server/pkg/models"
)

func (app *application) getYoutubeVideos(channel *models.Channel) ([]*models.Video, error) {
	return nil, nil
}

func (app *application) getYoutubeChannelData(login string, id string) (models.Channel, error) {
	var yc models.Channel
	url := "https://www.googleapis.com/youtube/v3/channels" +
		"?key=AIzaSyA2vHJcCFGgAKJv-g_l81lcNHxic9V4s3Y" +
		"&part=id,snippet" +
		"&fields=items(id,snippet(title,customUrl,thumbnails(default)))"

	if login != "" {
		url += "&forUsername=" + login
	} else if id != "" {
		url += "&id=" + id
	} else {
		return yc, errors.New("need to specify either login or id")
	}

	resp, err := app.httpClient.Get(url)
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
