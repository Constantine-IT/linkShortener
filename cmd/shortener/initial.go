package main

import (
	"database/sql"
	//	github.com/jackc/pgx/stdlib - драйвер PostgreSQL для доступа к БД с использованием пакета database/sql
	//	если хотим работать с БД напрямую, без database/sql надо использовать пакет - github.com/jackc/pgx/v4
	_ "github.com/jackc/pgx/stdlib"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
)

// initial - функция конструктор, инициализирующая хранилище URL и интерфейсы работы с файлом, хранящим URL
func initial(DatabaseDSN, FileStorage string) (strg storage.Datasource, err error) {
	//	Приоритетность в использовании ресурсов хранения информации URL (по убыванию приоритета):
	//	1.	Внешняя база данных, параметры соединения с которой задаются через DATABASE_DSN
	//	2.	Если БД не задана, то используем файловое хранилище (задаваемое через FILE_STORAGE_PATH) и оперативную память
	//	3.	Если не заданы ни БД, ни файловое хранилище, то работаем только с оперативной памятью - структура storage.Storage

	if DatabaseDSN != "" { //	если задана переменная среды DATABASE_DSN
		//	открываем connect с базой данных PostgreSQL 10+
		db, err := sql.Open("pgx", DatabaseDSN)
		if err != nil { //	при ошибке открытия, прерываем работу приложения
			return nil, err
		}
		//	тестируем доступность базы данных
		if err := db.Ping(); err != nil { //	если база недоступна, прерываем работу приложения
			return nil, err
		} else { //	если база данных доступна - создаём в ней структуры хранения

			//	готовим SQL-statement для создание таблицы shorten_urls, если её не существует
			stmt := `create table if not exists "shorten_urls" (
						"hash" text constraint hash_pk primary key not null,
   						"longurl" text constraint unique_longurl unique not null,
   						"userid" text not null)`
			_, err := db.Exec(stmt)
			if err != nil {	//	при ошибке в создании структур хранения в базе данных, прерываем работу приложения
				return nil, err
			}
		}
		//	если всё прошло успешно, возвращаем в качестве источника данных - базу данных
		strg = &storage.Database{DB: db}
	} else { //	если база данных не указана или недоступна
		//	возвращаем в качестве источника данных - структуру в оперативной памяти
		s := storage.Storage{Data: make(map[string]storage.RowStorage)}
		strg = &s

		//	опционально подключаем файл-хранилище URL
		if FileStorage != "" { //	если задан FILE_STORAGE_PATH
			//	порождаем reader и writer для файла-хранилища URL
			storage.URLreader, err = storage.NewReader(FileStorage)
			if err != nil {
				return nil, err
			}
			storage.URLwriter, err = storage.NewWriter(FileStorage)
			if err != nil {
				return nil, err
			}
			//	производим первичное заполнение хранилища URL в оперативной памяти из файла-хранилища URL
			err := storage.InitialURLFulfilment(&s)
			if err != nil {
				return nil, err
			}
		}
	}
	return strg, nil //	возращаем выбранный источник данных для URL
}
