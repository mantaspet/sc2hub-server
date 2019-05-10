package main

import (
	"bytes"
	"encoding/json"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"github.com/mantaspet/sc2hub-server/pkg/models/mock"
	"net/http"
	"testing"
)

func TestGetChannelsByCategory(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.router())
	defer ts.Close()

	twitchArr, err := json.Marshal([]models.Channel{*mock.TwitchChannel})
	if err != nil {
		t.Fatal(err)
	}

	emptyArr, err := json.Marshal([]models.Channel{})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid category", "/event-categories/1/channels", http.StatusOK, twitchArr},
		{"Invalid category", "/event-categories/3/channels", http.StatusOK, emptyArr},
		{"Valid category, valid platform", "/event-categories/1/channels?platform_id=1", http.StatusOK, twitchArr},
		{"Valid category, invalid platform", "/event-categories/1/channels?platform_id=2", http.StatusOK, emptyArr},
		{"Invalid category, invalid platform", "/event-categories/3/channels", http.StatusOK, emptyArr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			var test []models.Channel
			_ = json.Unmarshal(body, &test)
			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q; got %q", tt.wantBody, body)
			}
		})
	}
}

func TestGetTwitchAppAccessToken(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.router())
	defer ts.Close()

	code, _, body := ts.get(t, "/twitch/app-access-token")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	want := []byte("mockToken")
	if !bytes.Contains(body, want) {
		t.Errorf("want body to contain %q", want)
	}
}
