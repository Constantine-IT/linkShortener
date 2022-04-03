package handlers

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/models"
)

//	Базовый адрес для <shorten_URL>. Может быть переопределен в main.go
var Addr = "http://127.0.0.1:8080"

//	вспомогательная функция, создающая HASH из longURL, и сохраняющая связку HASH<==>URL в БД
//	возвращающает короткий URL для отправки клиенту
func saveShortURLlongURL(longURL string) string {

	// изготавливаем HASH из входящего URL с помощью MD5 hash algorithm
	md5URL := md5.Sum([]byte(longURL))
	hashURL := fmt.Sprintf("%X", md5URL[0:4])

	// вызов метода-вставки в структуру хранения связки HASH<==>longURL
	models.Insert(hashURL, longURL)

	// Изготавливаем  shortURL из адреса нашего сервера и HASH
	shortURL := strings.Join([]string{Addr, hashURL}, "/")
	return shortURL
}

//	Обработчики маршрутизатора

//	Обработчик POST с URL в виде JSON
func CreateShortURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	jsonURL, err := io.ReadAll(reader) // считываем JSON из тела запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//	описываем структуру JSON в запросе - {"url":"<some_url>"}
	type jsonURLBody struct {
		URL string `json:"url"`
	}
	//	создаеём экземпляр структуры для заполнения из JSON
	JSONBody := jsonURLBody{}

	//	парсим JSON и записываем результат в экземпляр структуры
	err = json.Unmarshal(jsonURL, &JSONBody)
	//	проверяем успешно ли парсится JSON
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//	Проверяем на пустую строку вместо URL
	if len(JSONBody.URL) == 0 {
		http.Error(w, "There is no URL in your request BODY!", http.StatusBadRequest)
		return
	}
	/*
		//	проверяем URL на недопустимые символы
		if strings.ContainsAny(JSONBody.URL, " !,*\n") {
			http.Error(w, "There are forbidden symbols in the URL!", http.StatusBadRequest)
			return
		}
	*/

	//	изготавливаем shortURL и сохраняем в БД связку HASH<==>URL
	shortURL := saveShortURLlongURL(JSONBody.URL)

	//	описываем структура создаваемого JSON вида {"result":"<shorten_url>"}
	type ResultURL struct {
		Result string `json:"result"`
	}
	//	создаем экземпляр структуры и вставляем в него короткий URL для отправки в JSON
	resultURL := ResultURL{
		Result: shortURL,
	}
	//	изготавливаем JSON
	shortJSONURL, err := json.Marshal(resultURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в JSON виде
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(shortJSONURL) //	пишем JSON с URL в тело ответа
}

//	Обработчик POST с URL в виде текста
func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	inURL, err := io.ReadAll(reader)
	//	проверяем на пустое тело запроса и/или другие ошибки чтения
	if err != nil || len(inURL) == 0 {
		http.Error(w, "There is no URL in your request BODY!", http.StatusBadRequest)
		return
	}

	longURL := string(inURL)
	/*
		//	проверяем URL на недопустимые символы
		if strings.ContainsAny(longURL, " !,*\n") {
			http.Error(w, "There are forbidden symbols in the URL!", http.StatusBadRequest)
			return
		}
	*/
	//	изготавливаем shortURL и сохраняем в базу связку HASH<==>longURL
	shortURL := saveShortURLlongURL(longURL)

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в виде текста
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL)) //	пишем URL в текстовом виде в тело ответа
}

//	Обработчик GET на адрес короткого URL
func GetShortURLHandler(w http.ResponseWriter, r *http.Request) {

	hashURL := chi.URLParam(r, "hashURL")

	//	проверяем указан ли HASH в коротком URL
	if hashURL == "" {
		http.Error(w, "ShortURL param is missed", http.StatusBadRequest)
		return
	}

	// Находим в базе URL соответствующий запрошенному HASH
	longURL, flag := models.Get(hashURL)
	if !flag {
		http.Error(w, "There is no such URL in our base!", http.StatusNotFound)
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя URL в заголовок в поле "location" и делая Redirect на него
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
