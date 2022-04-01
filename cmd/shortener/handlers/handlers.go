package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/models"
)

//	если задана переменная среды BASE_URL, то используем её как адрес для сокращенного URL
//	если не задана, то значение по умолчанию задается здесь в http://127.0.0.1:8080
var Addr = "http://127.0.0.1:8080"

//	вспомогательная общая функция, создающая HASH из longURL,
//	сохраняющая связку HASH<==>URL в базу данных, и возвращающая короткий URL для отправки клиенту

func saveShortURLlongURL(longURL string) string {

	// изготавливаем HASH из входящего URL с помощью MD5 hash algorithm
	md5URL := md5.Sum([]byte(longURL))
	hashURL := fmt.Sprintf("%X", md5URL[0:4])

	// вызов метода-вставки в структуру хранения связки HASH<==>longURL
	models.Insert(hashURL, longURL)

	// Изготавливаем  shortURL из адреса нашего сервера и HASH
	shortURL := strings.Join([]string{Addr, hashURL}, "/")
	return shortURL
}

//	Обработчики маршрутизатора
//	Обработчик POST с URL в виде JSON
func CreateShortURLJSONHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	jsonURL, err := io.ReadAll(r.Body) // считываем JSON из тела запроса
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type jsonURLBody struct { //	описываем структуру JSON в запросе - {"url":"<some_url>"}
		URL string `json:"url"`
	}
	JSONBody := jsonURLBody{} //	создаеём экземпляр структуры для заполнения из JSON

	err = json.Unmarshal(jsonURL, &JSONBody) //	парсим JSON и записываем результат в экземпляр структуры
	if err != nil {                          //	проверяем успешно ли парсится JSON
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(JSONBody.URL) == 0 { //	Проверяем на пустую строку вместо URL
		http.Error(w, "There is no URL in your request BODY!", http.StatusBadRequest)
		return
	}
	if strings.ContainsAny(JSONBody.URL, " !,*\n") { //	проверяем URL на недопустимые символы
		http.Error(w, "There are forbidden symbols in the URL!", http.StatusBadRequest)
		return
	}

	shortURL := saveShortURLlongURL(JSONBody.URL) //	изготавливаем shortURL и сохраняем в базу связку HASH<==>URL

	type ResultURL struct { //	описываем структура создаваемого JSON
		Result string `json:"result"`
	}
	resultURL := ResultURL{ //	создаем экземпляр структуры и вставляем в него короткий URL для отправки в JSON
		Result: shortURL,
	}
	shortJSONURL, err := json.Marshal(resultURL) //	изготавливаем JSON вида {"result":"<shorten_url>"}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в JSON виде
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(shortJSONURL) //	пишем JSON с URL в тело ответа
}

//	Обработчик POST с URL в виде текста
func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	inURL, err := io.ReadAll(r.Body)
	if err != nil || len(inURL) == 0 { //	проверяем на пустое тело запроса и/или другие ошибки чтения
		http.Error(w, "There is no URL in your request BODY!", http.StatusBadRequest)
		return
	}

	longURL := string(inURL)

	if strings.ContainsAny(longURL, " !,*\n") { //	проверяем URL на недопустимые символы
		http.Error(w, "There are forbidden symbols in the URL!", http.StatusBadRequest)
		return
	}

	shortURL := saveShortURLlongURL(longURL) //	изготавливаем shortURL и сохраняем в базу связку HASH<==>longURL

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в виде текста
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL)) //	пишем URL в текстовом виде в тело ответа
}

//	Обработчик GET на адрес короткого URL
func GetShortURLHandler(w http.ResponseWriter, r *http.Request) {

	hashURL := chi.URLParam(r, "hashURL")

	if hashURL == "" { //	проверяем указан ли HASH в коротком URL
		http.Error(w, "ShortURL param is missed", http.StatusBadRequest)
		return
	}

	longURL, flag := models.Get(hashURL) // Находим в базе URL соответствующий запрошенному HASH
	if !flag {
		http.Error(w, "There is no such URL in our base!", http.StatusNotFound)
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя URL в заголовок в поле "location" и делая Redirect на него
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	//	w.WriteHeader(201) //  Это для тестов в POSTMAN. На бою закомментировать.
}
