package handlers

import (
	"github.com/Constantine-IT/linkShortener/cmd/shortener/models"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Addr = "127.0.0.1:8080"

/* Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.  */

// Обработчик маршрутизатора

func CreateShortURLHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	inURL, err := io.ReadAll(r.Body)
	longURL := string(inURL)
	if err != nil || longURL == "" {
		http.Error(w, "There is no URL in your request BODY!", http.StatusBadRequest)
		return
	}
	if strings.ContainsAny(longURL, " !,*\n") {
		http.Error(w, "There are forbidden symbols in the URL!", http.StatusBadRequest)
		return
	}

	// пока в качестве ID короткого URL берем и сохраняем текущее локальное время в наносекундах
	// потом перепишем на использование HASH функции от аргумента - входящего длинного URL
	hashURL := strconv.FormatInt(time.Now().UnixNano(), 10)

	// для короткого URL храним только HASH - так меньше места занимает и быстрее поиск идет
	// вызов метода-вставки в структуру хранения связки HASH_ID<==>URL
	err = models.Insert(hashURL, longURL)
	if err != nil {
		http.Error(w, "Can't save URL to the database!", http.StatusInternalServerError)
		return
	}
	// Изготавливаем короткий URL из адреса нашего сервера и HASH
	shortURL := strings.Join([]string{"http:/", Addr, hashURL}, "/")

	// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в виде текста
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func GetShortURLHandler(w http.ResponseWriter, r *http.Request) {

	hashURL := chi.URLParam(r, "hashURL")

	if hashURL == "" {
		http.Error(w, "ShortURL param is missed", http.StatusBadRequest)
		return
	}

	longURL, err := models.Get(hashURL) // Находим в базе URL соответствующий запрошенному HASH
	if err != nil || longURL == "" {
		http.Error(w, "There is no such URL in our base!", http.StatusNotFound)
		return
	}

	// Изготавливаем и возвращаем ответ, вставляя URL в заголовок в поле location и делая редирект на него
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	//	w.WriteHeader(201) //  Это для тестов в POSTMAN. На бою закомментировать.
}
