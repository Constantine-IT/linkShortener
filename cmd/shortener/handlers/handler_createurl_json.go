package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
)

//	CreateShortURLJSONHandler - обработчик POST с URL в виде JSON
func (app *Application) CreateShortURLJSONHandler(w http.ResponseWriter, r *http.Request) {
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

	jsonURL, err := io.ReadAll(r.Body) // считываем JSON из тела запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.ErrorLog.Println("JSON body read error:" + err.Error())
		return
	}

	//	описываем структуру JSON в запросе - {"url":"<some_url>"}
	type jsonURLBody struct {
		URL string `json:"url"`
	}
	//	создаеём экземпляр структуры для заполнения из JSON
	jsonBody := jsonURLBody{}

	//	парсим JSON и записываем результат в экземпляр структуры
	err = json.Unmarshal(jsonURL, &jsonBody)
	//	проверяем успешно ли парсится JSON
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.ErrorLog.Println("JSON body parsing error:" + err.Error())
		return
	}

	//	проверяем URL на допустимый синтаксис
	if _, err := url.ParseRequestURI(jsonBody.URL); err != nil {
		http.Error(w, "Error with URL parsing", http.StatusBadRequest)
		app.ErrorLog.Println("Error with URL parsing" + err.Error())
		return
	}

	//	изготавливаем shortURL и сохраняем в БД связку HASH<==>URL + UserID
	shortURL, err := app.saveURLtoDB(jsonBody.URL, requestUserID.Value)
	if errors.Is(err, storage.ErrConflictRecord) {
		responseStatus = http.StatusConflict
	} else {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			app.ErrorLog.Println("URL save error:" + err.Error())
			return
		}
	}

	//	описываем структуру JSON для отправки ответа - {"result":"<shorten_url>"}
	type ResultURL struct {
		Result string `json:"result"`
	}
	//	создаем экземпляр структуры и вставляем в него короткий URL для отправки в JSON
	resultURL := ResultURL{
		Result: shortURL,
	}
	//	изготавливаем JSON
	shortJSONURL, err := json.Marshal(resultURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		app.ErrorLog.Println(err.Error())
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в JSON виде
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(shortJSONURL) //	пишем JSON с URL в тело ответа
}
