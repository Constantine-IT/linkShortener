package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/handlers"
	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
)

func main() {

	//	Приоритеты настроек:
	//	1.	Переменные окружения - ENV
	//	2.	Значения, задаваемые флагами при запуске из консоли
	//	3.	Значения по умолчанию.

	//	Считываем флаги запуска из командной строки и задаём значения по умолчанию, если флаг при запуске не указан
	ServerAddress := flag.String("a", "127.0.0.1:8080", "SERVER_ADDRESS - адрес запуска HTTP-сервера")
	BaseURL := flag.String("b", "http://127.0.0.1:8080", "BASE_URL - базовый адрес результирующего сокращённого URL")
	DatabaseDSN := flag.String("d", "", "DATABASE_DSN - адрес подключения к БД (PostgreSQL)")
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
	if u, flg := os.LookupEnv("DATABASE_DSN"); flg {
		*DatabaseDSN = u
	}
	if u, flg := os.LookupEnv("FILE_STORAGE_PATH"); flg {
		*FileStorage = u
	}

	//	определяем источник данных для работы с URL
	datasource, err := initial(*DatabaseDSN, *FileStorage)
	if err != nil {
		errorLog.Fatal(err)
	}

	//	инициализируем контекст нашего приложения, для определения в дальнейшем путей логирования ошибок и
	//	информационных сообщений; базового адреса нашего сервера и используемых хранилищ для URL
	app := &handlers.Application{
		ErrorLog:   errorLog,
		InfoLog:    infoLog,
		BaseURL:    *BaseURL,
		Datasource: datasource,
	}

	//	закрываем reader и writer для файла-хранилища URL
	defer storage.URLreader.Close()
	defer storage.URLwriter.Close()
	//	запуск сервера
	app.InfoLog.Printf("Server started at address: %s", *ServerAddress)
	srv := &http.Server{
		Addr:     *ServerAddress,
		ErrorLog: app.ErrorLog,
		Handler:  app.Routes(),
	}
	app.ErrorLog.Fatal(srv.ListenAndServe())
}
