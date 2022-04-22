package storage

import (
	"database/sql"
	"errors"
)

//	Database - структура хранилища URL, обертывающая пул подключений к базе данных
type Database struct {
	DB *sql.DB
}

// Insert - Метод для сохранения связки HASH и (<original_URL> + UserID)
func (d *Database) Insert(hash, longURL, userID string) error {
	//	пустые значения URL или UserID к вставке в хранилище не допускаются
	if hash == "" || longURL == "" || userID == "" {
		return ErrEmptyNotAllowed
	}
	//	готовим SQL-statement для вставки в базу и запускаем его на исполнение
	stmt := `insert into "shorten_urls" ("hash", "userid", "longurl") values ($1, $2, $3)`
	_, err := d.DB.Exec(stmt, hash, userID, longURL)
	if err != nil {
		return err
	}
	return nil
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
func (d *Database) GetByLongURL(longURL string) (hash string, flg bool) {
	var h string

	stmt := `select "hash" from "shorten_urls" where "longurl" = $1`
	err := d.DB.QueryRow(stmt, longURL).Scan(&h)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false
	}
	return h, true
}

// GetByUserID - Метод для нахождения списка сохраненных пар HASH и <original_URL> по UserID
func (d *Database) GetByUserID(userID string) ([]HashURLrow, bool) {
	var hash, longurl string
	hashRows := make([]HashURLrow, 0)

	stmt := `select "hash", "longurl" from "shorten_urls" where "userid" = $1`
	rows, err := d.DB.Query(stmt, userID)
	if err != nil || rows.Err() != nil {
		return nil, false
	}
	defer rows.Close()
	//	перебираем все строки выборки, добавляя связки HASH и <original_URL> в исходящий слайс
	for rows.Next() {
		err := rows.Scan(&hash, &longurl)
		if err != nil {
			return nil, false
		}
		hashRows = append(hashRows, HashURLrow{hash, longurl})
	}
	return hashRows, true
}
