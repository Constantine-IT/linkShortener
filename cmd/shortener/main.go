package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	h "github.com/Constantine-IT/linkShortener/cmd/shortener/handlers"
)

func main() {
	//если не задан ServerAddress в server.cfg, то по умолчанию запускаем сервер на 127.0.0.1:8080
	h.Addr = "127.0.0.1:8080"
	//	чтение файла конфигурации сервера	
	config, err := os.ReadFile("server.cfg")
	if err == nil {
		log.Printf("Читаем файл server.cfg \n %s", config)
	} else {
		log.Println(err.Error())
	}
	//	парсинг считанной конфигурации
	_, err = fmt.Sscanf(string(config), "ServerAddress %s", &h.Addr)
	if err != nil {
		log.Println(err.Error())
	} 
	log.Printf("Сервер будет запущен по адресу: %s", h.Addr)

	//	запуск сервера
	srv := &http.Server{
		Addr:    h.Addr,
		Handler: h.Routes(),
	}
	log.Fatal(srv.ListenAndServe())
}
