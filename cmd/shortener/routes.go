package main

import (
	"github.com/Constantine-IT/linkShortener/cmd/shortener/handlers"
	"net/http"
)

func Routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.ShortUrlHandler)

	return mux
}
