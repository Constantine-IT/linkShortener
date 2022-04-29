package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

//	GetShortURLHandler - обработчик GET на адрес короткого URL, возращает начальный URL по его короткой версии
func (app *Application) GetShortURLHandler(w http.ResponseWriter, r *http.Request) {
	//	считываем HASH из PATH входящего запроса
	hash := chi.URLParam(r, "hashURL")

	//	проверяем указан ли HASH в коротком URL
	if hash == "" {
		http.Error(w, "ShortURL param is missed", http.StatusBadRequest)
		app.ErrorLog.Println("ShortURL param is missed")
		return
	}

	// Находим в базе <original_URL> соответствующий запрошенному HASH
	longURL, flg := app.Datasource.Get(hash)

	switch flg {
	case 0: //	если URL в базе не найден
		http.Error(w, "There is no such URL in our database!", http.StatusNotFound)
		app.ErrorLog.Println("There is no such URL in our database!")
		return
	case 1: //	если URL в базе найден и нет пометки, что он "удалён"
		// Изготавливаем и возвращаем ответ, вставляя URL в заголовок в поле "location" и делая Redirect на него
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case 2: //	если URL в базе найден и на нём стоит пометка, что он "удалён"
		// Изготавливаем и возвращаем ответ со статусом 410 Gone - адрес удалён из базы
		w.WriteHeader(http.StatusGone)

	}
}
