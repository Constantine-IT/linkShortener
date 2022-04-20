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

	//	инициализируем контекст нашего приложения, для определения в дальнейшем путей логирования ошибок и
	//	информационных сообщений; базового адреса нашего сервера и используемых хранилищ для URL
	app := &handlers.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		BaseURL:  *BaseURL,
	}

	//	Приоритетность в использовании ресурсов сохранения информации URL (по убыванию приоритета):
	//	1.	Внешняя база данных, параметры соединения с которой задаются через DATABASE_DSN
	//	2.	Если БД не задана, то используем файловое хранилище (задаваемое через FILE_STORAGE_PATH) и оперативную память
	//	3.	Если не заданы ни БД, ни файловое хранилище, то работаем только с оперативной памятью - структура storage.Storage

	if *DatabaseDSN != "" { //	если база данных доступна, то работаем только с ней
		//	открываем connect с базой данных PostgreSQL 10+ по указанному DATABASE_DSN
		db, err := sql.Open("pgx", *DatabaseDSN)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
		defer db.Close()
		if err := db.Ping(); err == nil {
			//	Создание таблицы shorten_urls, если её не существует
			_, err := db.Exec(`create table if not exists "shorten_urls" (
    "hash" text constraint hash_pk primary key not null,
    "longurl" text constraint unique_longurl unique not null,
    "userid" text not null)`)
			if err != nil {
				app.ErrorLog.Println("DATABASE structure creation - " + err.Error())
			}
		} else {
			app.ErrorLog.Println("DATABASE open - " + err.Error())
			os.Exit(1)
		}
		app.InfoLog.Println("DATABASE creation - SUCCESS")
		app.Datasource = &storage.Database{DB: db}
		app.InfoLog.Println("DataBase connection has been established: " + *DatabaseDSN)
		app.InfoLog.Println("Server works only with DB, without file or RAM storage")
	} else {
		app.InfoLog.Println("DataBase wasn't set")
		s := storage.Storage{Data: make(map[string]storage.RowStorage)}
		app.Datasource = &s
		//	Первичное заполнение хранилища URL в оперативной памяти из файла-хранилища, если задан FILE_STORAGE_PATH
		if *FileStorage != "" {
			app.InfoLog.Printf("File storage with saved URL was found: %s", *FileStorage)
			storage.FileStorage = *FileStorage
			storage.InitialURLFulfilment(&s)
			app.InfoLog.Println("Saved URLs were loaded in RAM")
		} else {
			//	если файловое хранилище не задано, то работаем только в оперативной памяти
			app.InfoLog.Println("FileStorage wasn't set")
		}
		app.InfoLog.Println("Server works with RAM storage")
	}

	//	запуск сервера
	app.InfoLog.Printf("Server started at address: %s", *ServerAddress)
	srv := &http.Server{
		Addr:     *ServerAddress,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}
	app.ErrorLog.Fatal(srv.ListenAndServe())
}
