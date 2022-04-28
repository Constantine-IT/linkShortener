package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

//	CreateShortURLBatchHandler - обработчик POST с пакетом URL в виде JSON
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

	//	описываем структуру JSON в запросе -  {
	//        "correlation_id": "<строковый идентификатор>",
	//        "original_url": "<URL для сокращения>"
	//    }
	type incomingList struct {
		ShortURL string `json:""`
	}
	app.InfoLog.Println(string(jsonURL))
	//	создаеём экземпляр структуры для заполнения из JSON запроса
	incomingURLlist := make([]incomingList, 0)

	//	парсим JSON из запроса и записываем результат в экземпляр структуры incomingURLlist
	err = json.Unmarshal(jsonURL, &incomingURLlist)
	//	проверяем успешно ли парсится JSON
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.ErrorLog.Println("JSON body parsing error:" + err.Error())
		return
	}

	//	прогоняем цикл по всем строкам входящего списка URL
	for i := range incomingURLlist {
		//	проверяем URL на допустимый синтаксис
		if _, err := url.ParseRequestURI(incomingURLlist[i].ShortURL); err != nil {
			http.Error(w, "Error with URL parsing", http.StatusBadRequest)
			app.ErrorLog.Println("Error with URL parsing" + err.Error())
			return
		}

		//	изготавливаем shortURL и сохраняем в БД связку строку (HASH + <original_URL> + UserID)
		_, err := app.saveURLtoDB(incomingURLlist[i].ShortURL, requestUserID.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			app.ErrorLog.Println("URL save error:" + err.Error())
			return
		}
	}

	// Изготавливаем и возвращаем ответ, вставляя список коротких URL в тело ответа в JSON виде
	w.WriteHeader(http.StatusAccepted)
}
