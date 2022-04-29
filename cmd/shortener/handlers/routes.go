package handlers

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
)

type Application struct {
	ErrorLog   *log.Logger        //	журнал ошибок
	InfoLog    *log.Logger        //	журнал информационных сообщений
	BaseURL    string             //	базоовый адрес сервера
	Datasource storage.Datasource //	источник данных для хранения URL
}

func (app *Application) Routes() chi.Router {

	// определяем роутер chi
	r := chi.NewRouter()

	// зададим middleware для поддержки компрессии тел запросов и ответов
	r.Use(middleware.Compress(1, `text/plain`, `application/json`))
	r.Use(middleware.AllowContentEncoding(`gzip`))
	r.Use(app.DecompressGZIP)
	//	добавим функциональность аутентификации пользователя через симметрично подписанную куку,
	//	содержащую уникальный идентификатор пользователя
	r.Use(app.AuthCookie)
	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//	Эндпоинт GET /api/user/urls считывает UserID из cookie запроса и выдаёт все URL, сохраненные этим пользователем.
	//	Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
	//	Эндпоинт GET /ping проверяет доступность базы данных, выдает ответ с кодом 200, если доступна, и 500 - если не доступна.
	//	Эндпоинт POST / принимает в теле запроса в виде текста строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
	//	Эндпоинт POST /api/shorten - аналогичен предыдущему, но принимает в теле запроса JSON-объект {"url":"<some_url>"} и возвращает в теле ответа JSON-объект {"result":"<shorten_url>"}.
	//	Эндпоинт POST /api/shorten/batch, принимает в теле запроса множество URL для сокращения в формате JSON и возвращает сокращенные URL в JSON формате
	//	Эндпоинт DELETE /api/user/urls, принимает задания на удаление списка ранее сформированных URL
	r.Route("/", func(r chi.Router) {
		r.Delete("/api/user/urls", app.DeleteURLByUserIDHandler)
		r.Get("/api/user/urls", app.GetURLByUserIDHandler)
		r.Get("/{hashURL}", app.GetShortURLHandler)
		r.Get("/ping", app.PingDataBaseHandler)
		r.Get("/", app.GetShortURLHandler)
		r.Post("/api/shorten/batch", app.CreateShortURLBatchHandler)
		r.Post("/api/shorten", app.CreateShortURLJSONHandler)
		r.Post("/", app.CreateShortURLHandler)
	})

	return r
}
