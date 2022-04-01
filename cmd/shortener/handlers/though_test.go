package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//	These are integration test with following flow:
//	receive URL_for_shorting in POST body; create a <shorten_URL> from it and send <shorten_URL> to the client inside BODY,
//	receive <shorten_URL> from client with GET method and response to it with initial URL in the field "location" in header

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
			name:               "Going through test #1: with initial JSON body POST",
			initialRequest:     "/api/shorten",
			initialRequestType: "POST",
			body:               `{"url":"http://tudzqakmoorcb.net/bflsgr36aqo4x6/mmktfboj8"}`,
			secondRequestType:  http.MethodGet,
			want: want{
				inBetweenStatusCode:  http.StatusCreated,
				inBetweenContentType: "application/json",
				finalStatusCode:      http.StatusTemporaryRedirect,
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

			type BodyURL struct { //	описываем структуру JSON содержимого нашего ответа на POST c URL в JSON-формате
				Result string `json:"result"` //	URL лежит в JSON в формате {"result":"<shorten_url>"}
			}
			Body := BodyURL{} //	создаем экземпляр структуры и считываем в него JSON содержимое BODY
			_ = json.Unmarshal([]byte(body), &Body)
			body = Body.Result //	переопределяем переменную body, записыывая в неё URL, взятый из поля "result"

			//	теперь в body лежит <shorten_URL>, но тестовый сервер принимает только PATH без SCHEME и HOST
			//	так что вырезаем из <shorten_URL> прописанный в глобальной переменной handlers.Addr - BASE_URL
			body = strings.TrimPrefix(body, Addr)

			// используем содержимое body, как адрес запроса; само тело запроса оставляем пустым, так как это GET
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
			name:               "Going through test #2: with initial text/plain body POST",
			initialRequest:     "/",
			initialRequestType: "POST",
			body:               "http://tudzqakmoorcb.net/bflsgr36aqo4x6/mmktfboj8",
			secondRequestType:  http.MethodGet,
			want: want{
				inBetweenStatusCode:  http.StatusCreated,
				inBetweenContentType: "text/plain",
				finalStatusCode:      http.StatusTemporaryRedirect,
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

			//	теперь в body лежит <shorten_URL>, но тестовый сервер принимает только PATH без SCHEME и HOST
			//	так что вырезаем из <shorten_URL> прописанный в глобальной переменной handlers.Addr - BASE_URL
			body = strings.TrimPrefix(body, Addr)

			// используем содержимое body, как адрес запроса; само тело запроса оставляем пустым, так как это GET
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

	//	изменяем базовые политики redirect для HTTP-client - в случае редиректа, отменяем его и выдаём последний response
	http.DefaultClient.CheckRedirect = func(req *http.Request, previous []*http.Request) error {
		if len(previous) != 0 { // если были предыдущие запросы
			return http.ErrUseLastResponse //	возвращаем response, полученный на них
		}
		return nil
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
