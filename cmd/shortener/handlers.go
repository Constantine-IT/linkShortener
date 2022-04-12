package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/models"
)

//	вспомогательная функция, создающая HASH из longURL, и сохраняющая связку HASH<==>URL в БД
//	возвращает короткий URL для отправки клиенту
func (app *application) saveURLtoDB(longURL string) (string, error) {

	// изготавливаем HASH из входящего URL с помощью MD5 hash algorithm
	md5URL := md5.Sum([]byte(longURL))
	hashURL := fmt.Sprintf("%X", md5URL[0:4])

	// вызов метода-вставки в структуру хранения связки HASH<==>longURL
	err := models.Insert(hashURL, longURL, app.fileStorage, app.storage)

	// Изготавливаем  shortURL из адреса нашего сервера и HASH
	shortURL := strings.Join([]string{app.baseURL, hashURL}, "/")
	return shortURL, err
}

//	Обработчики маршрутизатора

//	Обработчик POST с URL в виде JSON
func (app *application) CreateShortURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	jsonURL, err := io.ReadAll(r.Body) // считываем JSON из тела запроса
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
		log.Println("Ошибка парсинга JSON-тела входящего запроса:\n" + err.Error())
		return
	}

	//	проверяем URL на допустимый синтаксис
	if _, err := url.ParseRequestURI(JSONBody.URL); err != nil {
		http.Error(w, "Error with parsing your URL!", http.StatusBadRequest)
		log.Println("Ошибка парсинга присланного URL:\n" + err.Error())
		return
	}

	//	изготавливаем shortURL и сохраняем в БД связку HASH<==>URL
	shortURL, err := app.saveURLtoDB(JSONBody.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Ошибка сохранения URL:" + err.Error())
		return
	}

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
func (app *application) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	inURL, err := io.ReadAll(r.Body)
	//	проверяем на ошибки чтения
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	longURL := string(inURL)
	//	проверяем URL на допустимый синтаксис
	if _, err := url.ParseRequestURI(longURL); err != nil {
		http.Error(w, "Error with parsing your URL!", http.StatusBadRequest)
		log.Println("Ошибка парсинга присланного URL:\n" + err.Error())
		return
	}
	//	изготавливаем shortURL и сохраняем в базу связку HASH<==>longURL
	shortURL, err := app.saveURLtoDB(longURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Ошибка сохранения URL:" + err.Error())
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в виде текста
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL)) //	пишем URL в текстовом виде в тело ответа
}

//	Обработчик GET на адрес короткого URL
func (app *application) GetShortURLHandler(w http.ResponseWriter, r *http.Request) {

	hashURL := chi.URLParam(r, "hashURL")

	//	проверяем указан ли HASH в коротком URL
	if hashURL == "" {
		http.Error(w, "ShortURL param is missed", http.StatusBadRequest)
		return
	}

	// Находим в базе URL соответствующий запрошенному HASH
	longURL, flag := models.Get(hashURL, app.storage)
	if !flag {
		http.Error(w, "There is no such URL in our base!", http.StatusNotFound)
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя URL в заголовок в поле "location" и делая Redirect на него
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
