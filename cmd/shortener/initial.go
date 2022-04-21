package main

import (
	"database/sql"
	"log"
	"os"
	//	github.com/jackc/pgx/stdlib - драйвер PostgreSQL для доступа к БД с использованием пакета database/sql
	//	если хотим работать с БД напрямую, без database/sql надо использовать пакет - github.com/jackc/pgx/v4
	_ "github.com/jackc/pgx/stdlib"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
)

// initial - функция конструктор, инициализирующая хранилище URL и интерфейсы работы с файлом, хранящим URL
func initial(DatabaseDSN, FileStorage string) (strg storage.Datasource) {
	var err error

	//	Приоритетность в использовании ресурсов хранения информации URL (по убыванию приоритета):
	//	1.	Внешняя база данных, параметры соединения с которой задаются через DATABASE_DSN
	//	2.	Если БД не задана, то используем файловое хранилище (задаваемое через FILE_STORAGE_PATH) и оперативную память
	//	3.	Если не заданы ни БД, ни файловое хранилище, то работаем только с оперативной памятью - структура storage.Storage

	if DatabaseDSN != "" { //	если база данных доступна, то работаем только с ней
		//	открываем connect с базой данных PostgreSQL 10+ по указанному DATABASE_DSN
		db, err := sql.Open("pgx", DatabaseDSN)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
		defer db.Close()
		//	тестируем доступность базы данных
		if err := db.Ping(); err != nil {
			log.Println("DATABASE open - " + err.Error())
			os.Exit(1)
		} else {
			//	если база данных доступна - создаём в ней структуры хранения
			//	готовим SQL-statement длясоздание таблицы shorten_urls, если её не существует
			stmt := `create table if not exists "shorten_urls" (
						"hash" text constraint hash_pk primary key not null,
   						"longurl" text constraint unique_longurl unique not null,
   						"userid" text not null)`
			_, err := db.Exec(stmt)
			if err != nil {
				log.Println("DATABASE structure creation - " + err.Error())
			}
		}
		log.Println("DATABASE creation - SUCCESS")
		strg = &storage.Database{DB: db}
		log.Println("DataBase connection has been established: " + DatabaseDSN)
		log.Println("Server works only with DB, without file or RAM storage")
	} else {
		log.Println("DataBase wasn't set")

		s := storage.Storage{Data: make(map[string]storage.RowStorage)}
		strg = &s

		//	если задан FILE_STORAGE_PATH
		if FileStorage != "" {
			log.Printf("File storage with saved URL was found: %s", FileStorage)
			//	порождаем reader и writer для файла хранилища URL
			storage.URLreader, err = storage.NewReader(FileStorage)
			if err != nil {
				log.Fatal(err)
			}
			storage.URLwriter, err = storage.NewWriter(FileStorage)
			if err != nil {
				log.Fatal(err)
			}

			//	первичное заполнение хранилища URL в оперативной памяти из файла-хранилища
			storage.InitialURLFulfilment(&s)
			log.Println("Saved URLs were loaded in RAM")
		} else {
			//	если файловое хранилище не задано, то работаем только в оперативной памяти
			log.Println("FileStorage wasn't set")
		}
		log.Println("Server works with RAM storage")
	}
	return strg
}
