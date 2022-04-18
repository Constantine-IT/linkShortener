package storage

import (
	"database/sql"
	"log"
)

//	Database - оределяет тип, который обертывает пул подключения sql.DB
type Database struct {
	DB *sql.DB
}

// Методы работы с базой данных - хранилищем URL

// Insert - Метод для сохранения связки короткого и длинного URL + UserID.
func (d *Database) Insert(hash, longURL, userID string) error {
	stmt := `insert into "shorten_urls" ("hash", "userid", "longurl") values ($1, $2, $3)`
	_, err := d.DB.Exec(stmt, hash, userID, longURL)
	if err != nil {
		log.Println("INSERT failed")
		return err
	} else {
		log.Println("INSERT was successful")
		return nil
	}
}

// Get - Метод для нахождения длинного URL по HASH из БД сохраненных URL
func (d *Database) Get(hash string) (longURL, userID string, flg bool) {
	var url string
	var user string
	stmt := `select "longurl", "userid" from "shorten_urls" where "hash" = $1`
	err := d.DB.QueryRow(stmt, hash).Scan(&url, &user)
	if err != nil {
		log.Println("SELECT by hash failed" + err.Error())
		return "", "", false
	}
	log.Println("SELECT by hash was successful")
	return url, user, true
}

// GetByUserID - Метод для нахождения списка сохраненных пар <shorten_URL> и <original_URL> по UserID
func (d *Database) GetByUserID(userID string) ([]HashURLrow, bool) {
	hashRows := make([]HashURLrow, 0)

	stmt := `select "hash", "longurl" from "shorten_urls" where "userid" = $1`
	rows, err := d.DB.Query(stmt, userID)
	if err != nil || rows.Err() != nil {
		log.Println("SELECT by UserID failed")
		return nil, false
	}
	defer rows.Close()
	for rows.Next() {
		var hash string
		var longurl string
		err := rows.Scan(&hash, &longurl)
		if err != nil {
			log.Println("SELECT by UserID failed")
			return nil, false
		}
		hashRows = append(hashRows, HashURLrow{hash, longurl})
	}
	log.Println("SELECT by UserID was successful")
	return hashRows, true
}

func (d *Database) Create() error {
	//	Создание таблицы shorten_urls

	_, err := d.DB.Exec(`create table "shorten_urls" (
    "hash" text constraint hash_pk primary key not null,
    "userid" text not null,
    "longurl" text not null)`)
	log.Println("table shorten_urls created")
	if err != nil {
		return err
	}
	return nil
}

//	Добавление индекса для созданного столбца
//CREATE INDEX idx_snippets_created ON userid(created), longurl(created);
