package main

import (
	h "github.com/Constantine-IT/linkShortener/cmd/shortener/handlers"
	"log"
	"net/http"
)

func main() {
	// Addr := "127.0.0.1:8080"

	srv := &http.Server{
		Addr:    h.Addr,
		Handler: Routes(),
	}

	log.Printf("Запуск сервера на %s", h.Addr)
	log.Fatal(srv.ListenAndServe())
}
