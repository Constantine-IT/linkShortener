package storage

import (
	"database/sql"
	"errors"
	"log"
)

//	Database - оределяет тип, который обертывает пул подключения sql.DB
type Database struct {
	DB *sql.DB
}

//	ErrConflictRecord - ошибка связанная с конфликтом записей в базе данных URL, когда пытаемся вставить запись, уже существующую в БД
var ErrConflictRecord = errors.New("storage-database: URL-record already exist")

// Методы работы с базой данных - хранилищем URL

// Insert - Метод для сохранения связки HASH и (<original_URL> + UserID)
func (d *Database) Insert(hash, longURL, userID string) error {
	//	готовим SQL-statement для вставки в базу и запускаем его на исполнение
	stmt := `insert into "shorten_urls" ("hash", "userid", "longurl") values ($1, $2, $3)`
	_, err := d.DB.Exec(stmt, hash, userID, longURL)
	if err != nil {
		log.Println("New URL INSERT - FAILED")
		return err
	} else {
		log.Println("New URL INSERT - SUCCESS")
		return nil
	}
}

// Get - Метод для нахождения <original_URL> и UserID по HASH
func (d *Database) Get(hash string) (longURL, userID string, flg bool) {
	var url string
	var user string

	stmt := `select "longurl", "userid" from "shorten_urls" where "hash" = $1`
	err := d.DB.QueryRow(stmt, hash).Scan(&url, &user)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", false
	}
	return url, user, true
}

// GetByLongURL - Метод для нахождения HASH по <original_URL>
func (d *Database) GetByLongURL(longURL string) (string, bool) {
	var hash string

	stmt := `select "hash" from "shorten_urls" where "longurl" = $1`
	err := d.DB.QueryRow(stmt, longURL).Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false
	}
	return hash, true
}

// GetByUserID - Метод для нахождения списка сохраненных пар HASH и <original_URL> по UserID
func (d *Database) GetByUserID(userID string) ([]HashURLrow, bool) {
	var hash, longurl string
	hashRows := make([]HashURLrow, 0)

	stmt := `select "hash", "longurl" from "shorten_urls" where "userid" = $1`
	rows, err := d.DB.Query(stmt, userID)
	if err != nil || rows.Err() != nil {
		log.Println("SELECT by UserID - FAILED")
		return nil, false
	}
	defer rows.Close()
	//	перебираем все строки выборки, добавляя связки HASH и <original_URL> в исходящий слайс
	for rows.Next() {
		err := rows.Scan(&hash, &longurl)
		if err != nil {
			log.Println("SELECT by UserID - FAILED")
			return nil, false
		}
		hashRows = append(hashRows, HashURLrow{hash, longurl})
	}
	log.Println("SELECT by UserID - SUCCESS")
	return hashRows, true
}

//	Create - метод создания структур хранения в базе данных URL
func (d *Database) Create() error {
	//	Создание таблицы shorten_urls
	_, err := d.DB.Exec(`create table "shorten_urls" (
    "hash" text constraint hash_pk primary key not null,
    "longurl" text constraint unique_longurl unique not null,
    "userid" text not null)`)
	if err != nil {
		return err
	}
	log.Println("DATABASE creation - SUCCESS")
	return nil
}
