package main

import (
	"encoding/json"
	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
	"net/http"
	"strings"
)

//	GetURLByUserIDHandler - обработчик GET для получения списка URL созданных пользователем с UserID
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
