package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

//	CreateShortURLBatchHandler - обработчик POST с пакетом URL в виде JSON
func (app *application) CreateShortURLBatchHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//	считываем UserID из cookie запроса
	requestUserID, err := r.Cookie("userid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.errorLog.Println("There is no userid in request cookie:" + err.Error())
		return
	}

	// считываем JSON из тела запроса
	jsonURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.errorLog.Println("JSON body read error:" + err.Error())
		return
	}

	//	описываем структуру JSON в запросе
	type incomingList struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	//	описываем структуру JSON в ответе
	type outgoingList struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	//	создаеём экземпляр структуры для заполнения из JSON запроса
	incomingURLlist := make([]incomingList, 0)

	//	создаеём экземпляр структуры для JSON ответа
	outgoingURLlist := make([]outgoingList, 0)

	//	парсим JSON из запроса и записываем результат в экземпляр структуры
	err = json.Unmarshal(jsonURL, &incomingURLlist)
	//	проверяем успешно ли парсится JSON
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.errorLog.Println("JSON body parsing error:" + err.Error())
		return
	}

	//	прогоняем цикл по всем срокам входящего JSON списка URL
	for i := range incomingURLlist {
		//	проверяем URL на допустимый синтаксис
		if _, err := url.ParseRequestURI(incomingURLlist[i].OriginalURL); err != nil {
			http.Error(w, "Error with URL parsing", http.StatusBadRequest)
			app.errorLog.Println("Error with URL parsing" + err.Error())
			return
		}

		//	изготавливаем shortURL и сохраняем в БД связку HASH<==>URL + UserID
		shortURL, err := app.saveURLtoDB(incomingURLlist[i].OriginalURL, requestUserID.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			app.errorLog.Println("URL save error:" + err.Error())
			return
		}
		
		//	создаем исходящий JSON список коротких URL
		outgoingURLlist = append(outgoingURLlist, outgoingList{incomingURLlist[i].CorrelationID, shortURL})
	}

	//	изготавливаем JSON для ответа
	shortJSONURL, err := json.Marshal(outgoingURLlist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.errorLog.Println(err.Error())
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в JSON виде
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(shortJSONURL) //	пишем JSON с URL в тело ответа
}
