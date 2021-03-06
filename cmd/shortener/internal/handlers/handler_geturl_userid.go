package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

//	GetURLByUserIDHandler - обработчик GET для получения списка URL созданных пользователем с UserID
func (app *Application) GetURLByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	//	считываем UserID из cookie запроса
	requestUserID, err := r.Cookie("userid")
	if err != nil {
		http.Error(w, "Cookie UserID error"+err.Error(), http.StatusInternalServerError)
		app.ErrorLog.Println("There is no userid in request cookie:" + err.Error())
		return
	}

	// Находим в базе URLs принадлежащие пользователю с данным UserID
	// вызов метода для нахождения в структуре хранения всех пар HASH<==>URL связанных с указанным UserID
	slicedURL, flg := app.Datasource.GetByUserID(requestUserID.Value)
	if !flg {
		http.Error(w, "There is no URL from this user in database", http.StatusNoContent)
		app.InfoLog.Println("There is no URL from this user in our database")
		return
	}

	//	Добавляем к каждому HASH базовый адрес ASE_URL
	for i := range slicedURL {
		slicedURL[i].Hash = strings.Join([]string{app.BaseURL, slicedURL[i].Hash}, "/")
	}
	//	кодируем информацию в JSON
	slicedJSONURL, err := json.Marshal(slicedURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.ErrorLog.Println(err.Error())
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя список URLs пользователя с UserID в тело ответа в JSON виде
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(slicedJSONURL) //	пишем JSON с URL в тело ответа
}
