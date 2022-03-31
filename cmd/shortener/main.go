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
	srvAddr := h.Addr
	//	чтение файла конфигурации сервера
	config, err := os.ReadFile("server.cfg")
	if err == nil {
		log.Printf("Читаем файл server.cfg \n %s", config)
	} else {
		log.Println(err.Error())
	}
	//	парсинг считанной конфигурации
	_, err = fmt.Sscanf(string(config), "ServerAddress %s", &srvAddr)
	if err != nil {
		log.Println(err.Error())
	}
	//	считываем переменные окружения: адрес запуска HTTP-сервера - SERVER_ADDRESS
	//	и базовый адрес результирующего сокращённого URL - BASE_URL
	if u, flag := os.LookupEnv("SERVER_ADDRESS"); flag {
		srvAddr = u //	если SERVER_ADDRESS задан, то стартуем наш HTTP-сервер на этом адресе
	}

	h.Addr = srvAddr

	if u, flag := os.LookupEnv("BASE_URL"); flag {
		h.Addr = u //	если переменная среды BASE_URL задана, то используем её как адрес для сокращенного URL
	}

	log.Printf("Сервер будет запущен по адресу: %s", srvAddr)

	//	запуск сервера
	srv := &http.Server{
		Addr:    srvAddr,
		Handler: h.Routes(),
	}
	log.Fatal(srv.ListenAndServe())
}
