package storage

import (
	"database/sql"
)

//	DatabaseModel - оределяет тип, который обертывает пул подключения sql.DB
type DatabaseModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой записи в базе дынных.
//func (m *ShortenURLtModel) Insert(hash, longURL, userID string) (int, error) {}

// Get - Метод для возвращения данных URL по его HASH.
//func (m *ShortenURLtModel) Get(id int) (*storage.Storage, error) {}
