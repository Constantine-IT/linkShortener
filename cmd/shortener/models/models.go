package models

import (
	"errors"
	"io"
	"log"
)

// таблица для хранения URL
var URLTable = make(map[string]string)

//	если FilePath задан - при перезапуске сервера, список <shorten_URL> сохраняется в этом файле
//	если FilePath не задан, то храним URL только в оперативной памяти и теряем при перезапуске.
var FilePath = ""

// Методы работы с хранилищем URL

//	Метод первичного заполнения БД из файла сохраненных URL при старте сервера
func InitialFulfilmentURLDB() {
	//	создаем экземпляр reader для файла-хранилища HASH<==>URL
	fileReader, err := newReader(FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fileReader.close()
	log.Println("Обнаружен файл сохраненных URL. Начинаем считывание:")
	for {
		//	считываем записи по одной из файла-хранилища HASH<==>URL
		readedURL, err := fileReader.read()
		//	когда дойдем до конца файла - выхоодим из цикла чтения
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//	записываем список считанных URL в log
		log.Println(readedURL)
		//	добавляем связку HASH<==>URL в таблицу в RAM
		URLTable[readedURL.hashURL] = readedURL.longURL
	}
}

// Метод для сохранения в БД связки короткого и длинного URL.
func Insert(shortURL, longURL string) {
	//	Проверяем наличие <shorten_URL> в списке сохраненных URL
	//	если такой URL уже есть в базе, то повторную вставку не производим
	if _, ok := URLTable[shortURL]; !ok {
		URLTable[shortURL] = longURL
		//	если файл для хранения URL не задан, то храним список только в RAM в URLTable
		if FilePath != "" {
			//	создаем экземпляр структуры хранения связки HASH<==>URL
			shortenURL := shortenURL{
				hashURL: shortURL,
				longURL: longURL,
			}
			//	создаем экземпляр writer для файла-хранилища HASH<==>URL
			writtenURL, err := newWriter(FilePath)
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
}

// Метод для нахождения длинного URL по HASH от <shorten_URL> из БД сохраненных URL
func Get(shortURL string) (longURL string, flag bool) {
	longURL, ok := URLTable[shortURL]
	return longURL, ok
}
