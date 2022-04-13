package storage

import (
	"errors"
	"log"
)

// Методы работы с хранилищем URL

// Insert - Метод для сохранения в БД связки короткого и длинного URL.
func Insert(shortURL, longURL, filePath string, storage *Storage) error {
	//	пустые значения URL к вставке в хранилище не допускаются
	if shortURL == "" || longURL == "" {
		return errors.New("empty value is not allowed")
	}
	//	Проверяем наличие <shorten_URL> в списке сохраненных URL
	//	если такой URL уже есть в базе, то повторную вставку не производим
	storage.mu.Lock()
	defer storage.mu.Unlock()
	if _, ok := storage.data[shortURL]; !ok {
		storage.data[shortURL] = longURL
		//	если файл для хранения URL не задан, то храним список только в RAM
		if filePath != "" {
			//	создаем экземпляр структуры хранения связки HASH<==>URL
			shortenURL := shortenURL{
				HashURL: shortURL,
				LongURL: longURL,
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
	}
	return nil
}

// Get - Метод для нахождения длинного URL по HASH от <shorten_URL> из БД сохраненных URL
func Get(shortURL string, storage *Storage) (longURL string, flag bool) {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	longURL, ok := storage.data[shortURL]
	return longURL, ok
}
