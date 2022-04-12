package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	m "github.com/Constantine-IT/linkShortener/cmd/shortener/models"
)

//goland:noinspection ALL
func TestResponseWithErrors(t *testing.T) {

	type want struct {
		statusCode  int
		contentType string
		body        string
	}
	//goland:noinspection SpellCheckingInspection
	tests := []struct {
		name        string
		request     string
		requestType string
		body        string
		want        want
	}{
		{
			name:        "Test #1: Request with empty body (without URL to short)",
			request:     "/",
			requestType: http.MethodPost,
			body:        "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Error with parsing your URL!\n",
			},
		},
		{
			name:        "Test #2: Request with non-absolute URL",
			request:     "/",
			requestType: http.MethodPost,
			body:        `test.com/ahshshd-wew?`,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Error with parsing your URL!\n",
			},
		},
		{
			name:        "Test #3: Request with method that not allowed",
			request:     "/",
			requestType: http.MethodPatch,
			body:        "http://test.com/ahshshd",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "",
				body:        "",
			},
		},
		{
			name:        "Test #4: Request URL that doesn't exist in database",
			request:     "/111",
			requestType: http.MethodGet,
			body:        "",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				body:        "There is no such URL in our base!\n",
			},
		},
		{
			name:        "Test #5: Request URL with too long PATH",
			request:     "/111/1223",
			requestType: http.MethodGet,
			body:        "",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				body:        "404 page not found\n",
			},
		},
		{
			name:        "Test #6: Request URL without HASH",
			request:     "/",
			requestType: http.MethodGet,
			body:        "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "ShortURL param is missed\n",
			},
		},
		{
			name:        "Test #7: Request with empty URL in JSON body",
			request:     "/api/shorten",
			requestType: http.MethodPost,
			body:        `{"url":""}`,
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "Error with parsing your URL!\n",
			},
		},
	}

	app := &application{
		errorLog:    log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		infoLog:     log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		baseURL:     "http://127.0.0.1:8080",
		storage:     m.NewStorage(),
		fileStorage: "",
		//database: &mysql.dbModel{DB: db},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := app.Routes()
			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, body := testSimpleRequest(t, ts, tt.requestType, tt.request, tt.body)
			defer resp.Body.Close()
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.body, body)
		})
	}
}

func testSimpleRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}