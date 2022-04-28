package storage

import (
	"database/sql"
	"errors"
)

//	Database - структура хранилища URL, обертывающая пул подключений к базе данных
type Database struct {
	DB *sql.DB
}

//	db - рабочий экземпляр структуры Database
var db Database

// Insert - метод для сохранения связки HASH и (<original_URL> + UserID)
func (d *Database) Insert(hash, longURL, userID string) error {
	//	пустые значения URL или UserID к вставке в хранилище не допускаются
	if hash == "" || longURL == "" || userID == "" {
		return ErrEmptyNotAllowed
	}

	//	начинаем тразакцию
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //	при ошибке выполнения - откатываем транзакцию

	//	готовим SQL-statement для вставки в базу
	stmt, err := tx.Prepare(`insert into "shorten_urls" ("hash", "userid", "longurl", "deleted") values ($1, $2, $3, false)`)
	if err != nil {
		return err
	}

	//	 запускаем SQL-statement на исполнение
	if _, err := stmt.Exec(hash, userID, longURL); err != nil {
		return err
	}

	return tx.Commit() //	при успешном выполнении вставки - фиксируем транзакцию
}

// DeleteByHashes - метод для пометки записей в базе данных как удаленные по их HASH и UserID
func (d *Database) DeleteByHashes(hashes []string, userID string) error {
	//	начинаем тразакцию
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //	при ошибке выполнения - откатываем транзакцию

	//	готовим SQL-statement для обновления статуса удаленных строк в базе данных
	stmt, err := tx.Prepare(`update "shorten_urls" set "deleted"=true where "hash" = $1 and "userid" = $2`)
	if err != nil {
		return err
	}
	for _, hash := range hashes {
		//	 запускаем SQL-statement на исполнение
		if _, err := stmt.Exec(hash, userID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Get - метод для нахождения <original_URL> и UserID по HASH
func (d *Database) Get(hash string) (longURL, userID string, flg int) {
	var url string
	var user string
	var isDeleted bool

	stmt := `select "longurl", "userid", "deleted"from "shorten_urls" where "hash" = $1`
	err := d.DB.QueryRow(stmt, hash).Scan(&url, &user, &isDeleted)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", 0
	}
	if isDeleted {
		return url, user, 2
	}
	return url, user, 1
}

// GetByLongURL - метод для нахождения HASH по <original_URL>
func (d *Database) GetByLongURL(longURL string) (hash string, flg bool) {
	var h string

	stmt := `select "hash" from "shorten_urls" where "longurl" = $1 and "deleted"=false`
	err := d.DB.QueryRow(stmt, longURL).Scan(&h)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false
	}
	return h, true
}

// GetByUserID - метод для нахождения списка сохраненных пар HASH и <original_URL> по UserID
func (d *Database) GetByUserID(userID string) ([]HashURLrow, bool) {
	var hash, longurl string
	hashRows := make([]HashURLrow, 0)

	stmt := `select "hash", "longurl" from "shorten_urls" where "userid" = $1 and "deleted"=false`
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

//	Close - метод, закрывающий reader и writer для файла-хранилища URL, а также connect к базе данных
func (d *Database) Close() error {
	var err error

	err = fileReader.Close()
	if err != nil {
		return err
	}

	err = fileWriter.Close()
	if err != nil {
		return err
	}

	err = db.DB.Close()
	if err != nil {
		return err
	}

	return nil
}
