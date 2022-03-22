package main

import (
	"log"
	"net/http"
)

func main() {
	addr := "localhost:8080"

	srv := &http.Server{
		Addr:    addr,
		Handler: Routes(),
	}

	log.Printf("Запуск сервера на %s", addr)
	log.Fatal(srv.ListenAndServe())
}
