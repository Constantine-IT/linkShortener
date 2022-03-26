package main

import (
	"fmt"
	h "github.com/Constantine-IT/linkShortener/cmd/shortener/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
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
