package main

import (
	"encoding/json"
	"errors"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type YoutubeSearchResult struct {
	Id struct {
		VideoId string
	}
}

type YoutubeVideo struct {
	Id             string
	ContentDetails struct {
		Duration string
	}
	Statistics struct {
		ViewCount string
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

// Gets video data by video IDs returned from youtube /search endpoint
func (app *application) getYoutubeVideoData(videos []YoutubeSearchResult) ([]YoutubeVideo, error) {
	var url strings.Builder

	url.WriteString("https://www.googleapis.com/youtube/v3/videos" +
		"?key=" + flgYoutubeApiKey +
		"&part=snippet,contentDetails,statistics" +
		"&fields=items(id," +
		"snippet(publishedAt,title,liveBroadcastContent,thumbnails(medium(url)))," +
		"contentDetails(duration)," +
		"statistics(viewCount))" +
		"&id=")

	for _, v := range videos {
		url.WriteString(v.Id.VideoId + ",")
	}

	res, err := app.httpClient.Get(strings.TrimRight(url.String(), ","))
	if err != nil {
		app.errorLog.Println(err.Error())
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

	return data.Items, nil
}

func parseYoutubeVideoDuration(durationString string) int {
	// youtube duration looks like this: PT1H25M30S
	var durationInSeconds int
	duration, err := time.ParseDuration(strings.ToLower(durationString[2:]))

	if err != nil {
		durationInSeconds = 0
	} else {
		durationInSeconds = int(duration.Seconds())
	}
	return durationInSeconds
}

// Gets the latest Youtube video data from a given channel
func (app *application) getYoutubeVideos(channel *models.Channel) ([]*models.Video, error) {
	url := "https://www.googleapis.com/youtube/v3/search" +
		"?key=" + flgYoutubeApiKey +
		"&channelId=" + channel.ID +
		"&type=video" +
		"&part=id" +
		"&fields=items(id)" +
		"&order=date" +
		"&maxResults=50"

	res, err := app.httpClient.Get(url)
	if err != nil {
		app.errorLog.Println(err.Error())
		return nil, err
	}

	type Response struct {
		Items []YoutubeSearchResult
	}

	var data Response
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	// Since /search endpoint does not provide all wanted fields, a different endpoint must be queried.
	// Can't just query that endpoint though, because it does not accept channel ID as filter param.
	youtubeVideos, err := app.getYoutubeVideoData(data.Items)
	if err != nil {
		return nil, err
	}

	var videos []*models.Video
	for _, v := range youtubeVideos {

		if v.Snippet.LiveBroadcastContent != "none" {
			continue
		}
		match, _ := regexp.MatchString("^[[[:ascii:]]+$", v.Snippet.Title) // to exclude videos with non-english titles
		if match != true {
			continue
		}
		if !app.matchesPattern([]string{v.Snippet.Title}, channel.IncludePatterns, channel.ExcludePatterns) {
			continue
		}
		createdAt, err := time.Parse("2006-01-02T15:04:05Z", v.Snippet.PublishedAt)
		if err != nil {
			createdAt = time.Now()
		}

		viewCount, err := strconv.Atoi(v.Statistics.ViewCount)
		if err != nil {
			viewCount = 0
		}

		video := &models.Video{
			ID:              v.Id,
			PlatformID:      2,
			EventCategoryID: channel.EventCategoryID,
			ChannelID:       channel.ID,
			Title:           v.Snippet.Title,
			Duration:        parseYoutubeVideoDuration(v.ContentDetails.Duration),
			ThumbnailURL:    v.Snippet.Thumbnails.Medium.Url,
			ViewCount:       viewCount,
			CreatedAt:       createdAt,
		}
		videos = append(videos, video)
	}

	return videos, nil
}

// Gets updated video data of videos already saved in database.
// Returns only a subset of video data, for use with app.videos.UpdateMetadata.
func (app *application) getExistingYoutubeVideoData(videos []*models.Video) ([]*models.Video, error) {
	var url strings.Builder

	url.WriteString("https://www.googleapis.com/youtube/v3/videos" +
		"?key=" + flgYoutubeApiKey +
		"&part=snippet,contentDetails,statistics" +
		"&fields=items(id," +
		"snippet(title,thumbnails(medium(url)))," +
		"contentDetails(duration)," +
		"statistics(viewCount))" +
		"&id=")

	for _, v := range videos {
		url.WriteString(v.ID + ",")
	}

	res, err := app.httpClient.Get(strings.TrimRight(url.String(), ","))
	if err != nil {
		app.errorLog.Println(err.Error())
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

	var updatedVideos []*models.Video
	for _, v := range data.Items {
		viewCount, err := strconv.Atoi(v.Statistics.ViewCount)
		if err != nil {
			viewCount = 0
		}

		video := &models.Video{
			ID:           v.Id,
			Title:        v.Snippet.Title,
			ThumbnailURL: v.Snippet.Thumbnails.Medium.Url,
			Duration:     parseYoutubeVideoDuration(v.ContentDetails.Duration),
			ViewCount:    viewCount,
		}
		updatedVideos = append(updatedVideos, video)
	}

	return updatedVideos, nil
}

var getYoutubeChannelData = func(login string, id string, httpClient *http.Client) (models.Channel, error) {
	var yc models.Channel
	url := "https://www.googleapis.com/youtube/v3/channels" +
		"?key=" + flgYoutubeApiKey +
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
