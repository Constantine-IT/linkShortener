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
	database    *storage.Database
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

	//	открываем connect с базой данных PostgreSQL по указанному DATABASE_DSN
	db, err := sql.Open("pgx", *DatabaseDSN)
	if err != nil {
		errorLog.Println(err.Error())
	}
	defer db.Close()

	app := &application{
		errorLog:    errorLog,
		infoLog:     infoLog,
		baseURL:     *BaseURL,
		storage:     storage.NewStorage(),
		database:    &storage.Database{DB: db},
		fileStorage: *FileStorage,
	}

	//	Приоритетность в использовании ресурсов сохранения информации URL (по убыванию приоритета):
	//	1.	Внешняя база данных, параметры соединения с которой задаются через DATABASE_DSN
	//	2.	Если БД не задана, то используем файловое хранилище (задаваемое через FILE_STORAGE_PATH) и оперативную память
	//	3.	Если не заданы ни БД, ни файловое хранилище, то работаем только с оперативной памятью - структура storage.Storage

	//	проверяем доступность базы данных
	if err := app.database.DB.Ping(); err == nil {
		//	если база данных доступна, то работаем только с ней
		app.database.Create()
		app.storage = nil
		app.fileStorage = ""
		infoLog.Println("DataBase connection has established: " + *DatabaseDSN)
		infoLog.Println("Server works only with DB, without file storages or RAM structures")
	} else {
		app.database = nil
		infoLog.Println("DataBase wasn't set")
		//	Первичное заполнение хранилища URL в оперативной памяти из файла-хранилища, если задан FILE_STORAGE_PATH
		if app.fileStorage != "" {
			infoLog.Printf("File storage with saved URL was found: %s", app.fileStorage)
			storage.InitialURLFulfilment(app.storage, app.fileStorage)
			infoLog.Println("Saved URLs were loaded in RAM")
			infoLog.Println("All new URLs will be saved in file storage")
		} else {
			//	если файловое хранилище не задано, то работаем только в оперативной памяти
			infoLog.Println("FileStorage wasn't set")
		}
		infoLog.Println("Server works with structures in RAM")
	}

	//	запуск сервера
	infoLog.Printf("Server started at address: %s", *ServerAddress)
	srv := &http.Server{
		Addr:     *ServerAddress,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}
	errorLog.Fatal(srv.ListenAndServe())
}
