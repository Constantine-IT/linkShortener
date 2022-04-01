package models

import (
	"log"
)

// Пока работаем с MAP, висящей в RAM, потом перепишем всё на работу с БД
var URLTable = make(map[string]string) // таблица для хранения URL
var FilePath = ""

// Методы работы с моделью данных

// Insert - Метод для вставки в базу связки короткого и длинного URL.
func Insert(shortURL, longURL string) {
	_, ok := URLTable[shortURL]
	if !ok {
		URLTable[shortURL] = longURL
		if FilePath != "" {
			shortenURL := ShortenURL{
				HashURL: shortURL,
				LongURL: longURL,
			}
			writtenURL, err := NewURLWriter(FilePath)
			if err != nil {
				log.Fatal(err)
			}
			defer writtenURL.Close()
			if err := writtenURL.WriteURL(&shortenURL); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Get - Метод для нахождения длинного URL по короткому URL.
func Get(shortURL string) (longURL string, flag bool) {
	longURL, ok := URLTable[shortURL]
	return longURL, ok
}
