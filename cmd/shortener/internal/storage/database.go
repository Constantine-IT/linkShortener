package storage

import (
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/lib/pq"
)

//	Database - структура хранилища URL, обертывающая пул подключений к базе данных
type Database struct {
	DB *sql.DB
}

// Insert - метод для сохранения связки HASH и (<original_URL> + UserID)
func (d *Database) Insert(hash, longURL, userID string) error {
	//	пустые значения URL или UserID к вставке в хранилище не допускаются
	if hash == "" || longURL == "" || userID == "" {
		return ErrEmptyNotAllowed
	}

	//	начинаем тразакцию
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //	при ошибке выполнения - откатываем транзакцию

	//	готовим SQL-statement для вставки в базу
	stmt, err := tx.Prepare(`insert into "shorten_urls" ("hash", "userid", "longurl", "deleted") values ($1, $2, $3, false)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	//	 запускаем SQL-statement на исполнение
	if _, err := stmt.Exec(hash, userID, longURL); err != nil {
		return err
	}

	return tx.Commit() //	при успешном выполнении вставки - фиксируем транзакцию
}

// Delete - метод помечает записи в базе данных, как удаленные по их HASH и UserID
func (d *Database) Delete(hashes []string, userID string) error {
	//	начинаем тразакцию
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() //	при ошибке выполнения - откатываем транзакцию

	//	готовим SQL-statement для обновления статуса удаленных строк в базе данных
	//	обновляем через BATCH UPDATE - вставляя в STATEMENT сразу срез из HASH
	stmt, err := tx.Prepare(`update "shorten_urls" set "deleted"=true where "hash" = any ($1) and "userid" = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	//	запускаем SQL-statement на исполнение передавая в него параметрами список HASH на удаление и UserID
	//	метка pq.Array нужна, чтобы драйвер понял, что в statement через подстановочную переменную передают массив
	if _, err := stmt.Exec(pq.Array(hashes), userID); err != nil {
		return err
	}

	return tx.Commit() //	при успешном выполнении вставки - фиксируем транзакцию
}

// Get - метод для нахождения <original_URL> по HASH
func (d *Database) Get(hash string) (longURL string, flg int) {
	var isDeleted bool

	stmt := `select "longurl", "deleted"from "shorten_urls" where "hash" = $1`
	err := d.DB.QueryRow(stmt, hash).Scan(&longURL, &isDeleted)
	if errors.Is(err, sql.ErrNoRows) {
		return "", 0
	} //	если HASH не найден, возвращаем flag=0
	if isDeleted {
		return "", 2
	} //	если HASH найден с пометкой "удалён", возвращаем flag=2
	return longURL, 1 //	если HASH найден и пометки "удалён" нет, возвращаем flag=1
}

// GetByLongURL - метод для нахождения HASH по <original_URL>
func (d *Database) GetByLongURL(longURL string) (hash string, flg bool) {

	stmt := `select "hash" from "shorten_urls" where "longurl" = $1 and "deleted"=false`
	err := d.DB.QueryRow(stmt, longURL).Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false
	}
	return hash, true
}

// GetByUserID - метод для нахождения списка сохраненных пар HASH и <original_URL> по UserID
func (d *Database) GetByUserID(userID string) ([]HashURLrow, bool) {
	var hash, longurl string
	hashRows := make([]HashURLrow, 0)

	stmt := `select "hash", "longurl" from "shorten_urls" where "userid" = $1 and "deleted" = false`
	rows, err := d.DB.Query(stmt, userID)
	if err != nil || rows.Err() != nil {
		return nil, false
	}
	defer rows.Close()
	//	перебираем все строки выборки, добавляя связки HASH и <original_URL> в исходящий срез
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
	//	при остановке сервера connect к базе данных
	err := d.DB.Close()
	if err != nil {
		return err
	}

	return nil
}
