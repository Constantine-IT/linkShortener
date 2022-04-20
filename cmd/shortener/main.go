package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	//	github.com/jackc/pgx/stdlib - драйвер PostgreSQL для доступа к БД с использованием пакета database/sql
	//	если хотим работать с БД напрямую, без database/sql надо использовать пакет - github.com/jackc/pgx/v4
	_ "github.com/jackc/pgx/stdlib"

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

	//	открываем connect с базой данных PostgreSQL 10+ по указанному DATABASE_DSN
	db, err := sql.Open("pgx", *DatabaseDSN)
	if err != nil {
		errorLog.Println(err.Error())
	}
	defer db.Close()

	//	инициализируем контекст нашего приложения, для определения в дальнейшем путей логирования ошибок и
	//	информационных сообщений; базового адреса нашего сервера и используемых хранилищ для URL
	app := &handlers.Application{
		ErrorLog:    errorLog,
		InfoLog:     infoLog,
		BaseURL:     *BaseURL,
		Storage:     storage.NewStorage(),
		Database:    &storage.Database{DB: db},
		FileStorage: *FileStorage,
	}

	//	Приоритетность в использовании ресурсов сохранения информации URL (по убыванию приоритета):
	//	1.	Внешняя база данных, параметры соединения с которой задаются через DATABASE_DSN
	//	2.	Если БД не задана, то используем файловое хранилище (задаваемое через FILE_STORAGE_PATH) и оперативную память
	//	3.	Если не заданы ни БД, ни файловое хранилище, то работаем только с оперативной памятью - структура storage.Storage

	//	проверяем доступность базы данных
	if err := app.Database.DB.Ping(); err == nil {
		//	если база данных доступна, то работаем только с ней
		err := app.Database.Create() //	создаем структуры хранения данных в БД
		if err != nil {
			app.ErrorLog.Println("DATABASE creation - " + err.Error())
		}
		app.InfoLog.Println("DATABASE creation - SUCCESS")
		app.Storage = nil
		app.FileStorage = ""
		infoLog.Println("DataBase connection has been established: " + *DatabaseDSN)
		infoLog.Println("Server works only with DB, without file or RAM storage")
	} else {
		app.Database = nil
		infoLog.Println("DataBase wasn't set")
		//	Первичное заполнение хранилища URL в оперативной памяти из файла-хранилища, если задан FILE_STORAGE_PATH
		if app.FileStorage != "" {
			infoLog.Printf("File storage with saved URL was found: %s", app.FileStorage)
			storage.InitialURLFulfilment(app.Storage, app.FileStorage)
			infoLog.Println("Saved URLs were loaded in RAM\nAll new URLs will be saved in the file storage")
		} else {
			//	если файловое хранилище не задано, то работаем только в оперативной памяти
			infoLog.Println("FileStorage wasn't set")
		}
		infoLog.Println("Server works with RAM storage")
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
