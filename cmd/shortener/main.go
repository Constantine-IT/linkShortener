package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	baseURL     string
	storage     *storage.Storage
	fileStorage string
}

func main() {

	//	Приоритеты настроек:
	//	1.	Переменные окружения - ENV
	//	2.	Значения, задаваемые флагами при запуске из консоли
	//	3.	Значения по умолчанию.

	//	Считываем флаги запуска из командной строки и задаём значения по умолчанию, если флаг при запуске не указан
	ServerAddress := flag.String("a", "127.0.0.1:8080", "SERVER_ADDRESS - адрес запуска HTTP-сервера")
	BaseURL := flag.String("b", "http://127.0.0.1:8080", "BASE_URL - базовый адрес результирующего сокращённого URL")
	FileStorage := flag.String("f", "", "FILE_STORAGE_PATH - путь до файла с сокращёнными URL")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//	считываем переменные окружения, если они заданы - переопределяем соответствующие локальные переменные:
	if u, flg := os.LookupEnv("SERVER_ADDRESS"); flg {
		*ServerAddress = u
	}
	if u, flg := os.LookupEnv("BASE_URL"); flg {
		*BaseURL = u
	}
	if u, flg := os.LookupEnv("FILE_STORAGE_PATH"); flg {
		*FileStorage = u
	}

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		baseURL:     *BaseURL,
		storage:     storage.NewStorage(),
		fileStorage: *FileStorage,
	}
	//	Первичное заполнение хранилища URL в оперативной памяти из файла-хранилища, если задан FILE_STORAGE_PATH
	if app.fileStorage != "" {
		infoLog.Printf("Обнаружен файл сохраненных URL: %s", *FileStorage)
		storage.InitialFulfilmentURLDB(app.storage, app.fileStorage)
		infoLog.Println("Сохраненные URL успешно считаны в RAM")
	}

	//	запуск сервера
	infoLog.Printf("Сервер будет запущен по адресу: %s", *ServerAddress)
	srv := &http.Server{
		Addr:    *ServerAddress,
		ErrorLog: errorLog,
		Handler: app.Routes(),
	}
	errorLog.Fatal(srv.ListenAndServe())
}
