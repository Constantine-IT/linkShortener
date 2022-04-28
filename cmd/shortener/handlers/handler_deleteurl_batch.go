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

	//	создаеём slice для заполнения из JSON запроса - {[ "a", "b", "c", "d", ...]}
	var hashes []string

	//	парсим JSON из запроса и записываем результат в экземпляр структуры incomingURLlist
	err = json.Unmarshal(jsonURL, &hashes)
	//	проверяем успешно ли парсится JSON
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.ErrorLog.Println("JSON body parsing error:" + err.Error())
		return
	}

	//	отправляем список HASH для удаления из базы данных
	if err := app.Datasource.DeleteByHashes(hashes, requestUserID.Value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.ErrorLog.Println("URL delete error:" + err.Error())
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя список коротких URL в тело ответа в JSON виде
	w.WriteHeader(http.StatusAccepted)
}
