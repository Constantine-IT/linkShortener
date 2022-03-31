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

func CreateShortURL_JSONHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	jsonURL, err := io.ReadAll(r.Body) // считываем JSON из тела запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type JsonURLBody struct {
		Url string `json:"url"`
	}
	JsonURL := JsonURLBody{}

	err = json.Unmarshal(jsonURL, &JsonURL) //	парсим JSON и записываем результат в структуру
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(JsonURL.Url) == 0 {
		http.Error(w, "There is no URL in your request BODY!", http.StatusBadRequest)
		return
	}

	longURL := fmt.Sprintf("%s", JsonURL.Url) // 	изготавливаем символьную строку с URL, считанным из JSON

	if strings.ContainsAny(longURL, " !,*\n") {
		http.Error(w, "There are forbidden symbols in the URL!", http.StatusBadRequest)
		return
	}

	shortURL := saveShortURLlongURL(longURL) //	изготавливаем shortURL и сохраняем в базу связку HASH<==>longURL
	/*
		// изготавливаем HASH из входящего URL с помощью MD5 hash algorithm
		md5URL := md5.Sum([]byte(longURL))
		hashURL := fmt.Sprintf("%X", md5URL[0:4])

		// вызов метода-вставки в структуру хранения связки HASH<==>URL
		models.Insert(hashURL, longURL)

		// Изготавливаем короткий URL из адреса нашего сервера и HASH
		shortURL := strings.Join([]string{"http:/", Addr, hashURL}, "/")


	*/
	type ResultURL struct {
		Result string `json:"result"`
	}
	resultURL := ResultURL{
		Result: shortURL,
	}
	shortJsonURL, err := json.Marshal(resultURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в виде текста
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(shortJsonURL)
}

func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	inURL, err := io.ReadAll(r.Body)
	if err != nil || len(inURL) == 0 {
		http.Error(w, "There is no URL in your request BODY!", http.StatusBadRequest)
		return
	}

	longURL := string(inURL)

	if strings.ContainsAny(longURL, " !,*\n") {
		http.Error(w, "There are forbidden symbols in the URL!", http.StatusBadRequest)
		return
	}

	shortURL := saveShortURLlongURL(longURL) //	изготавливаем shortURL и сохраняем в базу связку HASH<==>longURL
	/*
		// изготавливаем HASH из входящего URL с помощью MD5 hash algorithm
		md5URL := md5.Sum([]byte(longURL))
		hashURL := fmt.Sprintf("%X", md5URL[0:4])

		// вызов метода-вставки в структуру хранения связки HASH<==>URL
		models.Insert(hashURL, longURL)

		// Изготавливаем короткий URL из адреса нашего сервера и HASH
		shortURL := strings.Join([]string{"http:/", Addr, hashURL}, "/")


	*/
	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в виде текста
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL)) //nolint:errcheck
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
