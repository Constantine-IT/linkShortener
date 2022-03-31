package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/models"
)

var Addr = "127.0.0.1:8080"

// Обработчики маршрутизатора

func saveShortURLlongURL(longURL string) string {

	// изготавливаем HASH из входящего URL с помощью MD5 hash algorithm
	md5URL := md5.Sum([]byte(longURL))
	hashURL := fmt.Sprintf("%X", md5URL[0:4])

	// вызов метода-вставки в структуру хранения связки HASH<==>longURL
	models.Insert(hashURL, longURL)

	// Изготавливаем  shortURL из адреса нашего сервера и HASH
	shortURL := strings.Join([]string{"http:/", Addr, hashURL}, "/")
	return shortURL
}

func CreateShortURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	jsonURL, err := io.ReadAll(r.Body) // считываем JSON из тела запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type jsonURLBody struct {
		URL string `json:"url"`
	}
	JSONBody := jsonURLBody{}

	err = json.Unmarshal(jsonURL, &JSONBody) //	парсим JSON и записываем результат в структуру
	if err != nil {                          //	проверяем парсится ли JSON
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(JSONBody.URL) == 0 { //	Проверяем на пустую строку вместо URL
		http.Error(w, "There is no URL in your request BODY!", http.StatusBadRequest)
		return
	}
	if strings.ContainsAny(JSONBody.URL, " !,*\n") { //	проверяем URL на недопустимые символы
		http.Error(w, "There are forbidden symbols in the URL!", http.StatusBadRequest)
		return
	}

	shortURL := saveShortURLlongURL(JSONBody.URL) //	изготавливаем shortURL и сохраняем в базу связку HASH<==>longURL

	type ResultURL struct {
		Result string `json:"result"`
	}
	resultURL := ResultURL{
		Result: shortURL,
	}
	shortJSONURL, err := json.Marshal(resultURL) //	изготавливаем JSON вида "result":"url"
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в JSON виде
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(shortJSONURL) //	пишем JSON с URL в тело ответа
}

func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	inURL, err := io.ReadAll(r.Body)
	if err != nil || len(inURL) == 0 { //	проверяем на пустое тело запроса
		http.Error(w, "There is no URL in your request BODY!", http.StatusBadRequest)
		return
	}

	longURL := string(inURL)

	if strings.ContainsAny(longURL, " !,*\n") { //	проверяем URL на недопустимые символы
		http.Error(w, "There are forbidden symbols in the URL!", http.StatusBadRequest)
		return
	}

	shortURL := saveShortURLlongURL(longURL) //	изготавливаем shortURL и сохраняем в базу связку HASH<==>longURL

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в виде текста
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL)) //	пишем URL в текстовом виде в тело ответа
}

func GetShortURLHandler(w http.ResponseWriter, r *http.Request) {

	hashURL := chi.URLParam(r, "hashURL")

	if hashURL == "" {
		http.Error(w, "ShortURL param is missed", http.StatusBadRequest)
		return
	}

	longURL, flag := models.Get(hashURL) // Находим в базе URL соответствующий запрошенному HASH
	if !flag {
		http.Error(w, "There is no such URL in our base!", http.StatusNotFound)
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя URL в заголовок в поле location и делая Redirect на него
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	//	w.WriteHeader(201) //  Это для тестов в POSTMAN. На бою закомментировать.
}
