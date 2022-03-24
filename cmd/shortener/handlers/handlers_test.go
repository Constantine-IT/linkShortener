package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
		secondRequestType  string
		body               io.Reader
		want               want
	}{
		{
			name:               "Going through test #1",
			initialRequest:     "/",
			initialRequestType: "POST",
			body:               strings.NewReader("http://tudzqakmoorcb.net/bflsgr36aqo4x6/mmktfboj8"),
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
			firstRequest := httptest.NewRequest(tt.initialRequestType, tt.initialRequest, tt.body)
			w1 := httptest.NewRecorder()
			h1 := http.HandlerFunc(ShortURLHandler)
			h1.ServeHTTP(w1, firstRequest)
			firstResult := w1.Result()

			assert.Equal(t, tt.want.inBetweenStatusCode, firstResult.StatusCode)
			assert.Equal(t, tt.want.inBetweenContentType, firstResult.Header.Get("Content-Type"))

			inURL, err := ioutil.ReadAll(firstResult.Body)
			shortURL := string(inURL)
			require.NoError(t, err)
			err = firstResult.Body.Close()
			require.NoError(t, err)

			secondRequest := httptest.NewRequest(tt.secondRequestType, shortURL, nil)
			w2 := httptest.NewRecorder()
			h2 := http.HandlerFunc(ShortURLHandler)
			h2.ServeHTTP(w2, secondRequest)
			secondResult := w2.Result()

			assert.Equal(t, tt.want.finalStatusCode, secondResult.StatusCode)
			assert.Equal(t, tt.want.location, secondResult.Header.Get("Location"))
		})
	}
}
