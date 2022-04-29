package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

//	DeleteURLByUserIDHandler - обработчик DELETE со списком HASH в виде JSON
func (app *Application) DeleteURLByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//	считываем UserID из cookie запроса
	requestUserID, err := r.Cookie("userid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.ErrorLog.Println("There is no userid in request cookie:" + err.Error())
		return
	}

	// считываем JSON из тела запроса
	jsonURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.ErrorLog.Println("JSON body read error:" + err.Error())
		return
	}

	//	создаеём срез для заполнения из JSON запроса - {[ "a", "b", "c", "d", ...]}
	var hashes []string

	//	парсим JSON из запроса и записываем результат в срез с HASH
	err = json.Unmarshal(jsonURL, &hashes)
	//	проверяем успешно ли парсится JSON
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.ErrorLog.Println("JSON body parsing error:" + err.Error())
		return
	}

	//	отправляем список HASH для пометки, как "удалённые" в базе данных URL
	if err := app.Datasource.Delete(hashes, requestUserID.Value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.ErrorLog.Println("URL delete error:" + err.Error())
		return
	}

	// сообщаем в Response, что задание на "удаление" списка URL принято к исполнению
	w.WriteHeader(http.StatusAccepted)
}
