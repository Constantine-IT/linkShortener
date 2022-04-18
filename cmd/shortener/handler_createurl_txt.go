package main

import (
	"io"
	"net/http"
	"net/url"
)

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
