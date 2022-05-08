package storage

import (
	"database/sql"
	//	github.com/jackc/pgx/stdlib - драйвер PostgreSQL для доступа к БД с использованием пакета database/sql
	//	если хотим работать с БД напрямую, без database/sql надо использовать пакет - github.com/jackc/pgx/v4
	_ "github.com/jackc/pgx/stdlib"
)

// NewDatasource - функция конструктор, инициализирующая хранилище URL и интерфейсы работы с файлом, хранящим URL
func NewDatasource(DatabaseDSN, FileStorage string) (strg Datasource, err error) {
	//	Приоритетность в использовании ресурсов хранения информации URL (по убыванию приоритета):
	//	1.	Внешняя база данных, параметры соединения с которой задаются через DATABASE_DSN
	//	2.	Если БД не задана, то используем файловое хранилище (задаваемое через FILE_STORAGE_PATH) и оперативную память
	//	3.	Если не заданы ни БД, ни файловое хранилище, то работаем только с оперативной памятью - структура storage.Storage

	if DatabaseDSN != "" { //	если задана переменная среды DATABASE_DSN
		var err error
		var d Database
		//	открываем connect с базой данных PostgreSQL 10+
		d.DB, err = sql.Open("pgx", DatabaseDSN)
		if err != nil { //	при ошибке открытия, прерываем работу конструктора
			return nil, err
		}
		//	тестируем доступность базы данных
		if err := d.DB.Ping(); err != nil { //	если база недоступна, прерываем работу конструктора
			return nil, err
		} else { //	если база данных доступна - создаём в ней структуры хранения

			//	готовим SQL-statement для создания таблицы shorten_urls, если её не существует
			stmt := `create table if not exists "shorten_urls" (
						"hash" text constraint hash_pk primary key not null,
   						"longurl" text constraint unique_longurl unique not null,
   						"userid" text not null,
                        "deleted" boolean not null)`
			_, err := d.DB.Exec(stmt)
			if err != nil { //	при ошибке в создании структур хранения в базе данных, прерываем работу конструктора
				return nil, err
			}
		}
		//	если всё прошло успешно, возвращаем в качестве источника данных - базу данных
		strg = &Database{DB: d.DB}
	} else { //	если база данных не указана, возвращаем в качестве источника данных - структуру в оперативной памяти
		s := Storage{Data: make(map[string]RowStorage)}
		strg = &s

		//	опционально подключаем файл-хранилище URL
		if FileStorage != "" { //	если задан FILE_STORAGE_PATH, порождаем reader и writer для файла-хранилища URL
			fileReader, err = NewReader(FileStorage)
			if err != nil { //	при ошибке создания reader, прерываем работу конструктора
				return nil, err
			}
			fileWriter, err = NewWriter(FileStorage)
			if err != nil { //	при ошибке создания writer, прерываем работу конструктора
				return nil, err
			}
			//	производим первичное заполнение хранилища URL в оперативной памяти из файла-хранилища URL
			err := InitialURLFulfilment(&s)
			if err != nil { //	при ошибке первичного заполнения хранилища URL, прерываем работу конструктора
				return nil, err
			}
		}
	}
	return strg, nil //	если всё прошло ОК, то возращаем выбранный источник данных для URL
}
