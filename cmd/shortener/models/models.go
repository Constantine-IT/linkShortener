package models

import (
	"log"
)

// Пока работаем с MAP, висящей в RAM, потом перепишем всё на работу с БД
var UrlTable = make(map[string]string) // таблица для хранения URL
var FilePath = ""

// Методы работы с моделью данных

// Insert - Метод для вставки в базу связки короткого и длинного URL.
func Insert(shortURL, longURL string) {
	_, ok := UrlTable[shortURL]
	if !ok {
		UrlTable[shortURL] = longURL
		if FilePath != "" {
			shortenURL := ShortenURL{
				HashURL: shortURL,
				LongURL: longURL,
			}
			writtenURL, err := NewURLWriter(FilePath)
			defer writtenURL.Close()
			if err != nil {
				log.Fatal(err)
			}
			if err := writtenURL.WriteURL(&shortenURL); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Get - Метод для нахождения длинного URL по короткому URL.
func Get(shortURL string) (longURL string, flag bool) {
	longURL, ok := UrlTable[shortURL]
	return longURL, ok
}
