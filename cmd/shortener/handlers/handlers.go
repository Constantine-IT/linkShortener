package handlers

import (
	"github.com/Constantine-IT/linkShortener/cmd/shortener/models"
	"io"
	"net/http"
	//"net/url"
	"strconv"
	"strings"
	"time"
)

var Addr = "127.0.0.1:8080"

/* Эндпоинт POST / принимает в теле запроса строку URL для сокращения
и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле. */
/* Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.  */

// Обработчик маршрутизатора

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET or POST")
		http.Error(w, "Метод запрещен!", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodPost {
		defer r.Body.Close()
		inURL, err := io.ReadAll(r.Body)
		if err != nil || inURL == nil {
			http.Error(w, "Некорректный URL, не могу прочесть!", http.StatusBadRequest)
			return
		}

		// пока в качестве ID короткого URL берем и сохраняем текущее локальное время в наносекундах
		// потом перепишем на использование HASH функции от аргумента - входящего длинного URL
		hashURL := strconv.FormatInt(time.Now().UnixNano(), 10)

		// длинный URL храним в исходном виде без изменений, как считали в теле запроса.
		// вызов метода-вставки в структуру хранения связки HASH_ID<==>URL
		err = models.Insert(hashURL, string(inURL))
		if err != nil {
			http.Error(w, "Не могу запомнить URL!", http.StatusInternalServerError)
			return
		}

		/*
			parsURL, err := url.Parse(string(inURL))
			if err != nil {
				http.Error(w, "Некорректный URL, не могу распарсить!", http.StatusBadRequest)
				return
			}
			shortURL := strings.Join([]string{parsURL.Scheme, parsURL.Host}, "://") // соединяем префикс и домен через разделитель ://
			shortURL = strings.Join([]string{shortURL, hashURL}, "/")               // через / добавляем HASH - это и есть короткий URL
		*/
		shortURL := strings.Join([]string{"http:/", Addr, hashURL}, "/")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL)) // вставляем короткий URL в тело ответа в виде текста
	}

	if r.Method == http.MethodGet {
		id := r.URL.RequestURI()
		if id == "/" {
			http.Error(w, "Вы ввели не полный URL!", http.StatusBadRequest)
			return
		}
		id = strings.Trim(id, "/") // Обрезаем ведущий опостроф в PATH URL
		longURL, err := models.Get(id)
		if err != nil {
			http.Error(w, "В нашей базе такого URL не найдено!", http.StatusNotFound)
			return
		}
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		//w.WriteHeader(201) //  Это для тестов. На бою закомментировать.
	}
}
