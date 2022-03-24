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

/* Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.  */

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
		longURL := string(inURL)
		if err != nil || longURL == "" {
			http.Error(w, "Отсутствует URL в теле запроса!", http.StatusBadRequest)
			return
		}
		if strings.ContainsAny(longURL, " !,*\n") {
			http.Error(w, "URL в теле запроса содержит недопустимые символы!", http.StatusBadRequest)
			return
		}

		// пока в качестве ID короткого URL берем и сохраняем текущее локальное время в наносекундах
		// потом перепишем на использование HASH функции от аргумента - входящего длинного URL
		hashURL := strconv.FormatInt(time.Now().UnixNano(), 10)

		// длинный URL храним в исходном виде без изменений, как считали в теле запроса.
		// для короткого URL храним только HASH - так меньше места занимает и быстрее поиск идет
		// вызов метода-вставки в структуру хранения связки HASH_ID<==>URL
		err = models.Insert(hashURL, longURL)
		if err != nil {
			http.Error(w, "Не могу запомнить URL!", http.StatusInternalServerError)
			return
		}
		// Изготавливаем короткий URL из адреса нашего сервера и HASH
		shortURL := strings.Join([]string{"http:/", Addr, hashURL}, "/")

		// Изготавливаем и возвращаем ответ, вставляя короткий URL в тело ответа в виде текста
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
	}

	if r.Method == http.MethodGet {
		id := r.URL.RequestURI()       // Вырезаем PATH из входящего адреса запроса
		id = strings.Trim(id, "/")     // Обрезаем ведущий SLASH в PATH
		longURL, err := models.Get(id) // Находим в базе URL соответствующий запрошенному HASH
		if err != nil || longURL == "" {
			http.Error(w, "В нашей базе такого URL не найдено!", http.StatusNotFound)
			return
		}

		// Изготавливаем и возвращаем ответ, вставляя URL в заголовок в поле location и делая редирект на него
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		//w.WriteHeader(201) //  Это для тестов в POSTMAN. На бою закомментировать.
	}
}
