package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

//	CreateShortURLJSONHandler - обработчик POST с URL в виде JSON
func (app *Application) CreateShortURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	//	считываем UserID из cookie запроса
	requestUserID, err := r.Cookie("userid")
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.ErrorLog.Println("URL save error:" + err.Error())
		return
	}

	//	описываем структуру создаваемого JSON вида {"result":"<shorten_url>"}
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
	w.WriteHeader(http.StatusCreated)
	w.Write(shortJSONURL) //	пишем JSON с URL в тело ответа
}
