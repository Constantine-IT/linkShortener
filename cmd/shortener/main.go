package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	baseURL     string
	storage     *storage.Storage
	database    *storage.DatabaseModel
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

	//	если заданы параметры соединения с базой данных PostgreSQL, то открываем connect
	var db *sql.DB
	if *DatabaseDSN != "" {
		db, err := openDB(*DatabaseDSN)
		if err != nil {
			db = nil
			errorLog.Println("Can't open DataBase:" + err.Error())
		}
		defer db.Close()
	}

	infoLog.Println("DB is opened")

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		baseURL:     *BaseURL,
		storage:     storage.NewStorage(),
		database:    &storage.DatabaseModel{DB: db},
		fileStorage: *FileStorage,
	}

	infoLog.Println("APP struct created")

	//	Первичное заполнение хранилища URL в оперативной памяти из файла-хранилища, если задан FILE_STORAGE_PATH
	if *FileStorage != "" {
		infoLog.Printf("Обнаружен файл сохраненных URL: %s", *FileStorage)
		storage.InitialFulfilmentURLDB(app.storage, app.fileStorage)
		infoLog.Println("Сохраненные URL успешно считаны в RAM")
	}

	//	запуск сервера
	infoLog.Printf("Сервер будет запущен по адресу: %s", *ServerAddress)
	srv := &http.Server{
		Addr:     *ServerAddress,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}
	errorLog.Fatal(srv.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	//	открываем базу данных PostgreSQL версии 10+
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("DB is opening")
	return db, nil
}
