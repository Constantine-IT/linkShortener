package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/internal/storage"
)

//	CreateShortURLHandler - обработчик POST с URL в виде текста
func (app *Application) CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//	в случае успешного создания короткого URL, возвращаем ответ со статусом 201,
	//	если такой URL уже есть в базе, то возвращаем его со статусом 409 - Conflict
	//	в переменной responseStatus - будем хранить статус ответа:
	responseStatus := http.StatusCreated

	requestUserID, err := r.Cookie("userid") //	считываем "userid" из cookie запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.ErrorLog.Println("There is no userid in request cookie:" + err.Error())
		return
	}

	inURL, err := io.ReadAll(r.Body)
	//	проверяем на ошибки чтения
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.ErrorLog.Println(err.Error())
		return
	}

	longURL := string(inURL)
	//	проверяем URL на допустимый синтаксис
	if _, err := url.ParseRequestURI(longURL); err != nil {
		http.Error(w, "Error with URL parsing", http.StatusBadRequest)
		app.ErrorLog.Println("Error with URL parsing:" + err.Error())
		return
	}

	//	изготавливаем shortURL и сохраняем в базу связку HASH<==>longURL
	shortURL, err := app.saveURLtoDB(longURL, requestUserID.Value)
	if errors.Is(err, storage.ErrConflictRecord) {
		responseStatus = http.StatusConflict
	} else {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			app.ErrorLog.Println("URL save error:" + err.Error())
			return
		}
	}

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в виде текста
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(responseStatus)
	w.Write([]byte(shortURL)) //	пишем URL в текстовом виде в тело ответа
}
