package storage

import (
	"errors"
	"log"
)

// Методы работы с хранилищем URL

// Insert - Метод для сохранения в БД связки короткого и длинного URL.
func Insert(shortURL, longURL, userID, filePath string, storage *Storage) error {
	//	пустые значения URL или UserID к вставке в хранилище не допускаются
	if shortURL == "" || longURL == "" || userID == "" {
		return errors.New("empty value is not allowed")
	}
	//	Проверяем наличие <shorten_URL> в списке сохраненных URL
	//	если такой URL уже есть в базе, то повторную вставку не производим
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	//if _, ok := storage.data[shortURL]; !ok {
	storage.data[shortURL] = rowStorage{longURL, userID}
	//	если файл для хранения URL не задан, то храним список только в RAM
	if filePath != "" {
		//	создаем экземпляр структуры хранения связки HASH<==>URL
		shortenURL := shortenURL{
			HashURL: shortURL,
			LongURL: longURL,
			UserID:  userID,
		}
		//	создаем экземпляр writer для файла-хранилища HASH<==>URL
		writtenURL, err := newWriter(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer writtenURL.close()
		//	производим запись в файл-хранилище связки HASH<==>URL
		if err := writtenURL.write(&shortenURL); err != nil {
			log.Fatal(err)
		}
	}
	//}
	return nil
}

// Get - Метод для нахождения длинного URL по HASH от <shorten_URL> из БД сохраненных URL
func Get(shortURL string, storage *Storage) (longURL string, userID string, flag bool) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	if _, ok := storage.data[shortURL]; !ok {
		return "", "", false
	}
	return storage.data[shortURL].longURL, storage.data[shortURL].userID, true
}

type HashURLrow struct {
	HashURL string `json:"short_url"`
	LongURL string `json:"original_url"`
}

// GetByUserID - Метод для нахождения спика сохраненных длинных URL по UserID
func GetByUserID(userID string, storage *Storage) ([]HashURLrow, bool) {

	hashRows := make([]HashURLrow, 0)

	// блокируем хранилище URL на время считывания информации
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	for shortURL, row := range storage.data {
		if row.userID == userID {
			hashRows = append(hashRows, HashURLrow{shortURL, row.longURL})
		}
	}
	if len(hashRows) == 0 {
		return nil, false
	}
	return hashRows, true
}
