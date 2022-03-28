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

	// группируем все запросы в одном месте
	r.Route("/", func(r chi.Router) {
		// GET /HASH
		r.Get("/{hashURL}", GetShortURLHandler)
		r.Get("/", GetShortURLHandler)
		// POST /
		r.Post("/", CreateShortURLHandler)
	})

	return r
}
