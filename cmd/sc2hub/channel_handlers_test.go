package main

import (
	"net/http"
	"testing"
)

func TestGetTwitchAppAccessToken(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.router())
	defer ts.Close()

	code, _, body := ts.get(t, "/twitch/app-access-token")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	t.Log(string(body))

	want := "\"" + app.twitchAccessToken + "\""
	if string(body) != want {
		t.Errorf("want %q; got %q", want, string(body))
	}
}
