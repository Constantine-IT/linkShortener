package storage

import "errors"

//	Datasource - интерфейс источника данных URL
//	может реализовываться базой данных (Database) или структурами в оперативной памяти (Storage)
type Datasource interface {
	Insert(hash, longURL, userID string) error
	Get(hash string) (longURL, userID string, flg bool)
	GetByLongURL(longURL string) (hash string, flg bool)
	GetByUserID(userID string) ([]HashURLrow, bool)
}

//	HashURLrow - структура для выдачи (HASH + <original_URL>) по запросу строк с одинаковым UserID
//	используется в методах - Storage.GetByUserID и Database.GetByUserID
type HashURLrow struct {
	Hash    string `json:"short_url"`
	LongURL string `json:"original_url"`
}

//	shortenURL - структура для чтения/записи информации в файле-хранилище URL в виде JSON
//	используется в методах Storage.Insert и InitialURLFulfilment
type shortenURL struct {
	HashURL string `json:"hash-url"`
	LongURL string `json:"long-url"`
	UserID  string `json:"user-id"`
}

//	RowStorage - структура записи в хранилище URL в оперативной памяти
//	используется для формирования структуры Storage и метода Storage.Insert
type RowStorage struct {
	longURL string
	userID  string
}

//	ErrConflictRecord - ошибка возникающая, когда пытаемся вставить в базу запись c уже существующим URL
var ErrConflictRecord = errors.New("storage-database: URL-record already exist")

//	ErrEmptyNotAllowed - ошибка возникающая при попытке вставить пустое значение в любое поле структуры хранения URL
//	используется в методе Storage.Insert
var ErrEmptyNotAllowed = errors.New("ram-storage: empty value is not allowed")
