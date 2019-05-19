package main

import (
	"bytes"
	"errors"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"github.com/mantaspet/sc2hub-server/pkg/models/mock"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication() *application {
	app := &application{
		errorLog:          log.New(ioutil.Discard, "", 0),
		infoLog:           log.New(ioutil.Discard, "", 0),
		channels:          &mock.ChannelModel{},
		twitchAccessToken: "mockToken",
	}

	getTwitchChannelDataByLogin = func(login string, httpClient *http.Client) (models.Channel, error) {
		if login != "starcraft" {
			return models.Channel{}, errors.New("channel does not exist")
		}
		return *mock.Channels[0], nil
	}

	getYoutubeChannelData = func(login string, id string, httpClient *http.Client) (models.Channel, error) {
		return *mock.Channels[1], nil
	}

	isAuthenticated = func(app *application, endpoint func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			endpoint(w, r)
		})
	}

	return app
}

// Define a custom testServer type which anonymously embeds a httptest.Server
// instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add the cookie jar to the client, so that response cookies are stored
	// and then sent with subsequent requests.
	ts.Client().Jar = jar

	// Disable redirect-following for the client. Essentially this function
	// is called after a 3xx response is received by the client, and returning
	// the http.ErrUseLastResponse error forces it to immediately return the
	// received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// Implement a get method on our custom testServer type. This makes a GET
// request to a given url path on the test server, and returns the response
// status code, headers and body.
func (ts *testServer) sendRequest(t *testing.T, method string, urlPath string, payload []byte) (int, http.Header, []byte) {
	var reqBody io.Reader
	if payload != nil {
		reqBody = bytes.NewBuffer(payload)
	}
	req, err := http.NewRequest(method, ts.URL+urlPath, reqBody)

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
