package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortURLJSONHandler(t *testing.T) {

	type want struct {
		inBetweenStatusCode  int
		inBetweenContentType string
		finalStatusCode      int
		location             string
	}
	tests := []struct {
		name               string
		initialRequest     string
		initialRequestType string
		body               string
		secondRequestType  string
		want               want
	}{
		{
			name: "Going through test #1 with URL in initial JSON POST",
			//	get the URL, create a short URL from it and send it to the client,
			//	then get short URL from client and response to him with initial URL
			initialRequest:     "/api/shorten",
			initialRequestType: "POST",
			//body:               "{\"url\":\"http://tudzqakmoorcb.net/bflsgr36aqo4x6/mmktfboj8\"}",
			body:              `{"url":"http://tudzqakmoorcb.net/bflsgr36aqo4x6/mmktfboj8"}`,
			secondRequestType: "GET",
			want: want{
				inBetweenStatusCode:  201,
				inBetweenContentType: "application/json",
				finalStatusCode:      307,
				location:             "http://tudzqakmoorcb.net/bflsgr36aqo4x6/mmktfboj8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Routes()
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, body := testRequest(t, ts, tt.initialRequestType, tt.initialRequest, tt.body)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.inBetweenStatusCode, resp.StatusCode)
			assert.Equal(t, tt.want.inBetweenContentType, resp.Header.Get("Content-Type"))

			type BodyURL struct {
				Result string `json:"result"`
			}
			Body := BodyURL{}
			_ = json.Unmarshal([]byte(body), &Body)
			body = Body.Result

			//	в BODY лежит короткий URL, но тестовый сервер принимает только PATH без SCHEME и IP-адреса
			body = strings.TrimPrefix(body, "http://")
			body = strings.TrimPrefix(body, Addr)

			resp, _ = testRequest(t, ts, tt.secondRequestType, body, "")
			defer resp.Body.Close()
			assert.Equal(t, tt.want.finalStatusCode, resp.StatusCode)
			assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
		})
	}
}

func TestShortURLHandler(t *testing.T) {

	type want struct {
		inBetweenStatusCode  int
		inBetweenContentType string
		finalStatusCode      int
		location             string
	}
	tests := []struct {
		name               string
		initialRequest     string
		initialRequestType string
		body               string
		secondRequestType  string
		want               want
	}{
		{
			name: "Going through test #2 with URL in initial text/plain POST",
			//	get the URL, create a short URL from it and send it to the client,
			//	then get short URL from client and response to him with initial URL
			initialRequest:     "/",
			initialRequestType: "POST",
			body:               "http://tudzqakmoorcb.net/bflsgr36aqo4x6/mmktfboj8",
			secondRequestType:  "GET",
			want: want{
				inBetweenStatusCode:  201,
				inBetweenContentType: "text/plain",
				finalStatusCode:      307,
				location:             "http://tudzqakmoorcb.net/bflsgr36aqo4x6/mmktfboj8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Routes()
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, body := testRequest(t, ts, tt.initialRequestType, tt.initialRequest, tt.body)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.inBetweenStatusCode, resp.StatusCode)
			assert.Equal(t, tt.want.inBetweenContentType, resp.Header.Get("Content-Type"))

			//	в BODY лежит короткий URL, но тестовый сервер принимает только PATH без SCHEME и IP-адреса
			body = strings.TrimPrefix(body, "http://")
			body = strings.TrimPrefix(body, Addr)

			resp, _ = testRequest(t, ts, tt.secondRequestType, body, "")
			defer resp.Body.Close()
			assert.Equal(t, tt.want.finalStatusCode, resp.StatusCode)
			assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	ErrUseLastResponse := errors.New("net/http: use last response")

	http.DefaultClient.CheckRedirect = func(req *http.Request, previous []*http.Request) error {
		if len(previous) != 0 { //	В случае редиректа, блокируем его и возвращаем последний response
			return ErrUseLastResponse
		}
		return nil
	}

	resp, _ := http.DefaultClient.Do(req)
	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
