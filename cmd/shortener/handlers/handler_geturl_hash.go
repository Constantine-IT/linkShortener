package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

//	GetShortURLHandler - обработчик GET на адрес короткого URL, возращает начальный URL по его короткой версии
func (app *Application) GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var flg = true
	var longURL = ""

	hash := chi.URLParam(r, "hashURL")

	//	проверяем указан ли HASH в коротком URL
	if hash == "" {
		http.Error(w, "ShortURL param is missed", http.StatusBadRequest)
		app.ErrorLog.Println("ShortURL param is missed")
		return
	}

	// Находим в базе <original_URL> соответствующий запрошенному HASH
	longURL, _, flg = app.Datasource.Get(hash)

	if !flg {
		http.Error(w, "There is no such URL in our database!", http.StatusNotFound)
		app.ErrorLog.Println("There is no such URL in our database!")
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя URL в заголовок в поле "location" и делая Redirect на него
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
