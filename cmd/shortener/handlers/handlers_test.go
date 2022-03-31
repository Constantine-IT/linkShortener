package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseWithErrors(t *testing.T) {

	type want struct {
		statusCode  int
		contentType string
		body        string
	}
	tests := []struct {
		name        string
		request     string
		requestType string
		body        string
		want        want
	}{
		{
			name:        "test #1: Request with empty body (without URL to short)",
			request:     "/",
			requestType: "POST",
			body:        "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "There is no URL in your request BODY!\n",
			},
		},
		{
			name:        "test #2: Request with forbidden symbols in URL",
			request:     "/",
			requestType: "POST",
			body:        "http://test.com/ahshshd*!",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "There are forbidden symbols in the URL!\n",
			},
		},
		{
			name:        "test #3: Request with method that not allowed",
			request:     "/",
			requestType: "PATCH",
			body:        "http://test.com/ahshshd",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "",
				body:        "",
			},
		},
		{
			name:        "test #4: Request URL that doesn't exist in database",
			request:     "/111",
			requestType: "GET",
			body:        "",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				body:        "There is no such URL in our base!\n",
			},
		},
		{
			name:        "test #5: Request URL with too long PATH",
			request:     "/111/1223",
			requestType: "GET",
			body:        "",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				body:        "404 page not found\n",
			},
		},
		{
			name:        "test #6: Request URL without HASH",
			request:     "/",
			requestType: "GET",
			body:        "",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "ShortURL param is missed\n",
			},
		},
		{
			name:        "test #7: Request with URL in JSON body",
			request:     "/api/shorten",
			requestType: "POST",
			body:        "{\"url\":\"\"}",
			want: want{
				statusCode:  http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
				body:        "There is no URL in your request BODY!\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Routes()
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
