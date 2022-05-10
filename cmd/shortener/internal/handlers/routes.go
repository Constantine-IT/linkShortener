package handlers

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/internal/storage"
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

	//	маршруты сервера и их обработчики
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
