package handlers

import (
	"fmt"
	"github.com/Constantine-IT/linkShortener/cmd/shortener/models"
	"io"
	"net/http"
	"strconv"
	"time"
)

/* Эндпоинт POST / принимает в теле запроса строку URL для сокращения
и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле. */
/* Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.  */

// Обработчик маршрутизатора

func ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET or POST")
		//w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Метод запрещен!", 405)
		return
	}

	if r.Method == http.MethodPost {
		inUrl, err := io.ReadAll(r.Body) // читаем тело запроса

		if err != nil || inUrl == nil {
			http.Error(w, "Некорректный URL, введите заново!", 400)
			return
		}
		// пока в качестве ID короткой ссылки берем текущее локальное время в наносекундах
		// потом перепишем на использование HASH функции от аргумента - входящего URL
		shortUrl := strconv.FormatInt(time.Now().UnixNano(), 10)
		longUrl := string(inUrl)

		//	вызов метода-вставки в структуру хранения связки ID<==>URL
		err = models.Insert(shortUrl, longUrl)
		if err != nil {
			http.Error(w, "Не могу запомнить URL!", 500)
			return
		}
		w.WriteHeader(201)
		fmt.Fprintln(w, shortUrl)
	}

	if r.Method == http.MethodGet {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Вы не указали короткую ссылку!", 400)
			return
		}
		// вызов метода-запроса, выдающего созраненный URL по его ID
		longUrl, err := models.Get(id)
		if err != nil {
			http.Error(w, "URLs: записи с таким ID не найдено!", 404)
			return
		}
		w.Header().Set("Location", longUrl)
		w.WriteHeader(307)
	}
	//	w.Write([]byte("Привет, это сервис изготовления коротких ссылок!"))
}
