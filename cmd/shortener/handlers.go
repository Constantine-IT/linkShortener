package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
)

//	saveURLtoDB - вспомогательная функция, создающая HASH из связки (URL + UserID),
//	сохраняет связку HASH<==>URL+UserID и возвращает короткий URL для отправки клиенту
func (app *application) saveURLtoDB(longURL, userID string) (string, error) {
	var err error
	// изготавливаем HASH из входящего URL с помощью MD5 hash algorithm
	mdSum := md5.Sum([]byte(longURL + userID))
	hash := fmt.Sprintf("%X", mdSum[0:4])

	// вызов метода-вставки в структуру хранения связки HASH<==>URL+UserID
	if app.storage != nil {
		err = app.storage.Insert(hash, longURL, userID, app.fileStorage)
	}
	if app.database != nil {
		err = app.database.Insert(hash, longURL, userID)
	}

	if err == nil {
		// Изготавливаем  <shorten_URL> из базового адреса нашего сервера и HASH
		shortURL := strings.Join([]string{app.baseURL, hash}, "/")
		return shortURL, nil
	} else {
		return "", err
	}
}

//	Обработчики маршрутизатора

//	CreateShortURLJSONHandler - обработчик POST с URL в виде JSON
func (app *application) CreateShortURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//	считываем UserID из cookie запроса
	requestUserID, err := r.Cookie("userid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.errorLog.Println("There is no userid in request cookie:" + err.Error())
		return
	}
	jsonURL, err := io.ReadAll(r.Body) // считываем JSON из тела запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.errorLog.Println("JSON body read error:" + err.Error())
		return
	}

	//	описываем структуру JSON в запросе - {"url":"<some_url>"}
	type jsonURLBody struct {
		URL string `json:"url"`
	}
	//	создаеём экземпляр структуры для заполнения из JSON
	jsonBody := jsonURLBody{}

	//	парсим JSON и записываем результат в экземпляр структуры
	err = json.Unmarshal(jsonURL, &jsonBody)
	//	проверяем успешно ли парсится JSON
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.errorLog.Println("JSON body parsing error:" + err.Error())
		return
	}

	//	проверяем URL на допустимый синтаксис
	if _, err := url.ParseRequestURI(jsonBody.URL); err != nil {
		http.Error(w, "Error with URL parsing", http.StatusBadRequest)
		app.errorLog.Println("Error with URL parsing" + err.Error())
		return
	}

	//	изготавливаем shortURL и сохраняем в БД связку HASH<==>URL + UserID
	shortURL, err := app.saveURLtoDB(jsonBody.URL, requestUserID.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.errorLog.Println("URL save error:" + err.Error())
		return
	}

	//	описываем структуру создаваемого JSON вида {"result":"<shorten_url>"}
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
		app.errorLog.Println(err.Error())
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в JSON виде
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(shortJSONURL) //	пишем JSON с URL в тело ответа
}

//	CreateShortURLHandler - обработчик POST с URL в виде текста
func (app *application) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//	считываем UserID из cookie запроса
	requestUserID, err := r.Cookie("userid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.errorLog.Println("There is no userid in request cookie:" + err.Error())
		return
	}

	inURL, err := io.ReadAll(r.Body)
	//	проверяем на ошибки чтения
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.errorLog.Println(err.Error())
		return
	}

	longURL := string(inURL)
	//	проверяем URL на допустимый синтаксис
	if _, err := url.ParseRequestURI(longURL); err != nil {
		http.Error(w, "Error with URL parsing", http.StatusBadRequest)
		app.errorLog.Println("Error with URL parsing:" + err.Error())
		return
	}
	//	изготавливаем shortURL и сохраняем в базу связку HASH<==>longURL
	shortURL, err := app.saveURLtoDB(longURL, requestUserID.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.errorLog.Println("Error with saving URL:" + err.Error())
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в виде текста
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL)) //	пишем URL в текстовом виде в тело ответа
}

//	Обработчик GET на адрес короткого URL
func (app *application) GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var flg = true
	var longURL = ""

	hash := chi.URLParam(r, "hashURL")

	//	проверяем указан ли HASH в коротком URL
	if hash == "" {
		http.Error(w, "ShortURL param is missed", http.StatusBadRequest)
		app.errorLog.Println("ShortURL param is missed")
		return
	}

	// Находим в базе URL соответствующий запрошенному HASH

	// вызов метода для нахождения в структуре хранения связки HASH<==>URL+UserID
	if app.storage != nil {
		longURL, _, flg = app.storage.Get(hash)
	}
	if app.database != nil {
		longURL, _, flg = app.database.Get(hash)
	}

	if !flg {
		http.Error(w, "There is no such URL in our base!", http.StatusNotFound)
		app.errorLog.Println("There is no such URL in our base!")
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя URL в заголовок в поле "location" и делая Redirect на него
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

//	GetURLByUserIDHandler - обработчик GET для получения списка URL созданных пользователем
func (app *application) GetURLByUserIDHandler(w http.ResponseWriter, r *http.Request) {

	var flg = true
	var slicedURL []storage.HashURLrow

	//	считываем UserID из cookie запроса
	requestUserID, err := r.Cookie("userid")
	if err != nil {
		http.Error(w, "Cookie UserID error"+err.Error(), http.StatusInternalServerError)
		app.errorLog.Println("There is no userid in request cookie:" + err.Error())
		return
	}

	// Находим в базе URLs принадлежащие пользователю с данным UserID
	// вызов метода для нахождения в структуре хранения всех пар HASH<==>URL связанных с указанным UserID
	if app.storage != nil {
		slicedURL, flg = app.storage.GetByUserID(requestUserID.Value)
	}
	if app.database != nil {
		slicedURL, flg = app.database.GetByUserID(requestUserID.Value)
	}

	if !flg {
		http.Error(w, "There is no URL from this user in database", http.StatusNoContent)
		app.errorLog.Println("There is no URL from this user in our database")
		return
	}

	//	Добавляем к каждому HASH базовый адрес ASE_URL
	for i := range slicedURL {
		slicedURL[i].Hash = strings.Join([]string{app.baseURL, slicedURL[i].Hash}, "/")
	}
	//	кодируем информацию в JSON
	slicedJSONURL, err := json.Marshal(slicedURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.errorLog.Println(err.Error())
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя список URLs пользователя с UserID в тело ответа в JSON виде
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(slicedJSONURL) //	пишем JSON с URL в тело ответа
}

//	PingDataBaseHandler - обработчик проверки доступности базы данных через GET /ping
func (app *application) PingDataBaseHandler(w http.ResponseWriter, r *http.Request) {
	if app.database == nil {
		http.Error(w, "DataBase environment variable wasn't set", http.StatusInternalServerError)
		return
	}
	err := app.database.DB.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Error(w, http.StatusText(200), http.StatusOK)
	}
}
