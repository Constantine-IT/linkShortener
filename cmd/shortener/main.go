package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	h "github.com/Constantine-IT/linkShortener/cmd/shortener/handlers"
	m "github.com/Constantine-IT/linkShortener/cmd/shortener/models"
)

func main() {

	//	Приоритеты настроек:
	//	1.	Переменные окружения - ENV
	//	2.	Значения, задаваемые флагами при запуске из консоли
	//	3.	Значения по умолчанию.

	//	Считываем флаги запуска из командной строки и задаём значения по умолчанию, если флаг при запуске не указан
	serverAddress := flag.String("a", "127.0.0.1:8080", "SERVER_ADDRESS - адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "http://127.0.0.1:8080", "BASE_URL - базовый адрес результирующего сокращённого URL")
	fileStoragePath := flag.String("f", "", "FILE_STORAGE_PATH - путь до файла с сокращёнными URL")
	flag.Parse()

	//	значание флагов записываем в локальные переменные
	srvAddr := *serverAddress
	h.Addr = *baseURL
	m.FilePath = *fileStoragePath

	//	считываем переменные окружения, если они заданы - переопределяем соответствующие локальные переменные:
	if u, flg := os.LookupEnv("SERVER_ADDRESS"); flg {
		srvAddr = u
	}
	if u, flg := os.LookupEnv("BASE_URL"); flg {
		h.Addr = u
	}
	if u, flg := os.LookupEnv("FILE_STORAGE_PATH"); flg {
		m.FilePath = u
	}

	//	Первичное заполнение БД <shorten_URL> из файла-хранилища, если задан FILE_STORAGE_PATH
	if m.FilePath != "" {
		m.InitialFulfilmentURLDB()
	}

	//	запуск сервера
	log.Printf("Сервер будет запущен по адресу: %s", srvAddr)
	srv := &http.Server{
		Addr:    srvAddr,
		Handler: h.Routes(),
	}
	log.Fatal(srv.ListenAndServe())
}
