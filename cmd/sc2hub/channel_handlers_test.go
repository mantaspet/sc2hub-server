package main

import (
	"bytes"
	"encoding/json"
	"github.com/mantaspet/sc2hub-server/pkg/models/mock"
	"net/http"
	"testing"
)

type channelHandlerTest struct {
	name     string
	url      string
	payload  []byte
	wantCode int
	wantBody []byte
}

func runChannelHandlerTests(t *testing.T, tests []channelHandlerTest, method string) {
	app := newTestApplication()
	ts := newTestServer(t, app.router())

	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.sendRequest(t, method, tt.url, tt.payload)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q; got %q", tt.wantBody, body)
			}
		})
	}
}

func TestGetChannelsByCategory(t *testing.T) {
	t.Parallel()

	wantBody1, err := json.Marshal(mock.GetCategoryChannels(2, 0))
	wantBody2, err := json.Marshal(mock.GetCategoryChannels(3, 0))
	wantBody3, err := json.Marshal(mock.GetCategoryChannels(1, 1))
	wantBody4, err := json.Marshal(mock.GetCategoryChannels(1, 2))
	wantBody5, err := json.Marshal(mock.GetCategoryChannels(3, 3))
	if err != nil {
		t.Fatal(err)
	}

	tests := []channelHandlerTest{
		{"Valid category", "/event-categories/2/channels",
			nil, http.StatusOK, wantBody1},
		{"Non-existing category", "/event-categories/3/channels",
			nil, http.StatusOK, wantBody2},
		{"Valid category, valid platform", "/event-categories/1/channels?platform_id=1",
			nil, http.StatusOK, wantBody3},
		{"Valid category, invalid platform", "/event-categories/1/channels?platform_id=2",
			nil, http.StatusOK, wantBody4},
		{"Non-existing category, non-existing platform", "/event-categories/3/channels?platform_id=3",
			nil, http.StatusOK, wantBody5},
		{"Invalid category", "/event-categories/asd/channels?platform_id=1",
			nil, http.StatusBadRequest, []byte("strconv.Atoi: parsing \"asd\": invalid syntax\n")},
		{"Invalid platform", "/event-categories/1/channels?platform_id=asd",
			nil, http.StatusBadRequest, []byte("strconv.Atoi: parsing \"asd\": invalid syntax\n")},
	}

	runChannelHandlerTests(t, tests, "GET")
}

func TestGetAllTwitchChannels(t *testing.T) {
	t.Parallel()

	wantBody, err := json.Marshal(mock.GetPlatformChannels(1))
	if err != nil {
		t.Fatal(err)
	}

	tests := []channelHandlerTest{
		{"All twitch", "/channels/twitch", nil, http.StatusOK, wantBody},
	}

	runChannelHandlerTests(t, tests, "GET")
}

func TestAddChannelToCategory(t *testing.T) {
	t.Parallel()

	type Request struct {
		URL string
	}

	payload1, err := json.Marshal(Request{"https://twitch.tv/starcraft/"})
	payload2, err := json.Marshal(Request{"https://www.youtube.com/channel/UCk3w4CQ_SlLH4V0-V6WjFZg/videos"})
	payload3, err := json.Marshal(Request{"https://www.youtube.com/user/WCSStarCraft/videos"})
	payload4, err := json.Marshal(Request{"https://google.com/"})
	payload5, err := json.Marshal(Request{"https://twitch.tv/starcraftss/"})
	want1, err := json.Marshal(mock.Channels[0])
	want2, err := json.Marshal(mock.Channels[1])
	want3, err := json.Marshal(map[string]string{"url": "Must be a valid twitch.tv or youtube.com channel URL"})
	want4, err := json.Marshal(map[string]string{"url": "Channel does not exist"})
	want5, err := json.Marshal(map[string]string{"url": "This channel is already in database"})
	if err != nil {
		t.Fatal(err)
	}

	tests := []channelHandlerTest{
		{"Valid twitch channel", "/event-categories/1/channels",
			payload1, http.StatusOK, want1},
		{"Valid youtube channel with ID", "/event-categories/1/channels",
			payload2, http.StatusOK, want2},
		{"Valid youtube channel with user login", "/event-categories/1/channels",
			payload3, http.StatusOK, want2},
		{"Invalid category ID", "/event-categories/-2/channels",
			payload3, http.StatusBadRequest, []byte("must specify a valid category ID")},
		{"Invalid channel URL", "/event-categories/1/channels",
			payload4, http.StatusUnprocessableEntity, want3},
		{"Non existing channel URL", "/event-categories/1/channels",
			payload5, http.StatusUnprocessableEntity, want4},
		{"Duplicate channel", "/event-categories/2/channels",
			payload2, http.StatusUnprocessableEntity, want5},
	}

	runChannelHandlerTests(t, tests, "POST")
}

func TestDeleteCategoryChannel(t *testing.T) {
	t.Parallel()

	tests := []channelHandlerTest{
		{"Invalid category ID", "/event-categories/-2/channels/42508152",
			nil, http.StatusBadRequest, []byte("must specify a valid category ID")},
		{"Non existing category", "/event-categories/3/channels/42508152",
			nil, http.StatusNotFound, []byte("no channel found in category")},
		{"Non existing channel", "/event-categories/2/channels/55555",
			nil, http.StatusNotFound, []byte("no channel found in category")},
		{"Existing channel", "/event-categories/1/channels/42508152",
			nil, http.StatusOK, []byte("channel was deleted")},
	}

	runChannelHandlerTests(t, tests, "DELETE")
}

func TestGetTwitchAppAccessToken(t *testing.T) {
	t.Parallel()

	app := newTestApplication()
	ts := newTestServer(t, app.router())
	defer ts.Close()

	code, _, body := ts.sendRequest(t, "GET", "/twitch/app-access-token", nil)

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	want := []byte("mockToken")
	if !bytes.Contains(body, want) {
		t.Errorf("want body to contain %q", want)
	}
}
