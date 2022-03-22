package handlers

import (
	"github.com/Constantine-IT/linkShortener/cmd/shortener/models"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/* Эндпоинт POST / принимает в теле запроса строку URL для сокращения
и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле. */
/* Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.  */

// Обработчик маршрутизатора

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
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
			http.Error(w, "Некорректный URL, введите заново!", http.StatusBadRequest)
			return
		}
		u, err := url.Parse(string(inURL))
		if err != nil {
			http.Error(w, "Некорректный URL, введите заново!", http.StatusBadRequest)
			return
		}
		hostURL := u.Host
		// пока в качестве ID короткой ссылки берем текущее локальное время в наносекундах
		// потом перепишем на использование HASH функции от аргумента - входящего URL
		hashURL := strconv.FormatInt(time.Now().UnixNano(), 10)
		shortURL := strings.Join([]string{hostURL, hashURL}, "/")
		// Длинный URL храним в исходном виде без изменений
		// короткий URL храним без префиксов, в виде домен/HASH

		//	вызов метода-вставки в структуру хранения связки ID<==>URL
		err = models.Insert(shortURL, string(inURL))
		if err != nil {
			http.Error(w, "Не могу запомнить URL!", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortURL))
		// fmt.Fprintln(w, shortURL)
	}

	if r.Method == http.MethodGet {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Вы не указали короткую ссылку!", http.StatusBadRequest)
			return
		}
		u, err := url.Parse(string(id))
		if err != nil {
			http.Error(w, "Некорректный URL, введите заново!", http.StatusBadRequest)
			return
		}
		hostURL := u.Host
		hashURL := u.Path
		// предварительно режем из короткого URL все префиксы и другую муть, так как мы его в таком виде и храним.
		// оставляем только домен и HASH
		shortURL := strings.Join([]string{hostURL, hashURL}, "")
		// вызов метода-запроса, выдающего созраненный URL по его сокращенному виду.
		longURL, err := models.Get(shortURL)
		if err != nil {
			http.Error(w, "URLs: записи с таким ID не найдено!", http.StatusNotFound)
			return
		}
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		// w.WriteHeader(http.StatusOK) //  Это для тестов. На бою закомментировать.
	}
}
