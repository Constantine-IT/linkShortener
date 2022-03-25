package main

import (
	h "github.com/Constantine-IT/linkShortener/cmd/shortener/handlers"
	"log"
	"net/http"
)

func main() {

	srv := &http.Server{
		Addr:    h.Addr,
		Handler: h.Routes(),
	}

	log.Printf("Запуск сервера на %s", h.Addr)
	log.Fatal(srv.ListenAndServe())
}
