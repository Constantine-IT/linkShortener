package models

import (
	"log"
)

// таблица для хранения URL
var URLTable = make(map[string]string)

//	если FilePath задан - при перезапуске сервера БД <shorten_URL> сохраняется в этом файле
//	если FilePath не задан, то храним БД URL только в оперативной памяти и теряем при перезапуске.
var FilePath = ""

// Методы работы с хранилищем URL

//	Метод первичного заполнения БД из файла сохраненных URL при старте сервера
func InitialFulfilmentURLDB() {
	fileReader, err := NewURLReader(FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fileReader.Close()
	log.Println("Из файла считаны сохраненные URL:")
	for {
		readedURL, err := fileReader.ReadURL()
		if readedURL == nil {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(readedURL)
		URLTable[readedURL.HashURL] = readedURL.LongURL
	}
}

// Метод для сохранения в БД связки короткого и длинного URL.
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

// Метод для нахождения длинного URL по HASH короткого URL из БД сохраненных URL
func Get(shortURL string) (longURL string, flag bool) {
	longURL, ok := URLTable[shortURL]
	return longURL, ok
}
