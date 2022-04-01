package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes() chi.Router {

	// определяем роутер chi
	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//	Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
	//	Эндпоинт POST /api/shorten, принимающий в теле запроса JSON-объект {"url":"<some_url>"} и возвращающий в ответ объект {"result":"<shorten_url>"}.
	//	Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
	r.Route("/", func(r chi.Router) {
		r.Get("/{hashURL}", GetShortURLHandler)
		r.Get("/", GetShortURLHandler)
		r.Post("/api/shorten", CreateShortURLJSONHandler)
		r.Post("/", CreateShortURLHandler)
	})

	return r
}
