package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	h "github.com/Constantine-IT/linkShortener/cmd/shortener/handlers"
)

func main() {
	h.Addr = "127.0.0.1:8080" //	адрес запуска HTTP-сервера. Значение по умолчанию.
	//	чтение файла конфигурации
	config, err := os.ReadFile("server.cfg")
	if err == nil {
		log.Printf("Читаем файл server.cfg \n %s", config)
	} else {
		log.Println(err.Error())
	}
	//	парсинг считанной конфигурации
	_, err = fmt.Sscanf(string(config), "ServerAddress %s", &h.Addr)
	if err == nil {
		log.Printf("Сервер будет запущен по адресу: %s", h.Addr)
	} else {
		log.Println(err.Error())
	}

	//	запуск сервера, если не задан адрес в server.cfg, по умолчанию запускаем сервер на 127.0.0.1:8080
	srv := &http.Server{
		Addr:    h.Addr,
		Handler: h.Routes(),
	}

	log.Printf("Запуск сервера на %s", h.Addr)
	log.Fatal(srv.ListenAndServe())
}
