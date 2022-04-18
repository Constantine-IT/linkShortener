package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

//	GetShortURLHandler - обработчик GET на адрес короткого URL, возращает начальный URL по его короткой версии
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
