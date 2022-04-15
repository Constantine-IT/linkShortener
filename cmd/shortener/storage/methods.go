package storage

import (
	"errors"
	"log"
)

// Методы работы с хранилищем URL

// Insert - Метод для сохранения в БД связки короткого и длинного URL + UserID.
func Insert(shortURL, longURL, userID, filePath string, storage *Storage) error {
	//	пустые значения URL или UserID к вставке в хранилище не допускаются
	if shortURL == "" || longURL == "" || userID == "" {
		return errors.New("empty value is not allowed")
	}

	//	Блокируем структуру храниения в опративной памяти на время записи информации
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	//	сохраняем URL в оперативной памяти в структуре storage.data
	//	кажда запись - это сопоставленная с HASH структура из (URL + UserID) - rowStorage
	storage.data[shortURL] = rowStorage{longURL, userID}
	//	если файл для хранения URL не задан, то храним список только в RAM
	if filePath != "" {
		//	создаем экземпляр структуры хранения связки HASH<==>URL+UserID
		shortenURL := shortenURL{
			HashURL: shortURL,
			LongURL: longURL,
			UserID:  userID,
		}
		//	создаем экземпляр writer для файла-хранилища
		writtenURL, err := newWriter(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer writtenURL.close()
		//	производим сохранение в файл связки HASH<==>URL+UserID
		if err := writtenURL.write(&shortenURL); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

// Get - Метод для нахождения длинного URL по HASH из БД сохраненных URL
func Get(shortURL string, storage *Storage) (longURL string, userID string, flag bool) {

	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	//	если записи с запрашиваемым HASH нет в базе, то выставялем FLAG в положение FALSE
	if _, ok := storage.data[shortURL]; !ok {
		return "", "", false
	}
	return storage.data[shortURL].longURL, storage.data[shortURL].userID, true
}

//	Структура записи связки HASH<==>URL для выдачи по запросу всех строк с одинаковым UserID
type HashURLrow struct {
	HashURL string `json:"short_url"`
	LongURL string `json:"original_url"`
}

// GetByUserID - Метод для нахождения списка сохраненных пар <shorten_URL> и <original_URL> по UserID
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

	//	если записей с таким UserID не найдено - выставляем FLAG в положение FALSE
	if len(hashRows) == 0 {
		return nil, false
	} else {
		//	если нашли, то возвращаем список пар <shorten_URL> и <original_URL> для запрашиваемого UserID
		return hashRows, true
	}
}
