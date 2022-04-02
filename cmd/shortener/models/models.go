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
	//	создаем экземпляр READER из файла-хранилища HASH<==>URL
	fileReader, err := NewURLReader(FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer fileReader.Close()
	log.Println("Из файла считаны сохраненные URL:")
	for {
		//	считываем записи по одной из файла-хранилища HASH<==>URL
		readedURL, err := fileReader.ReadURL()
		//	когда дойдем до конца файла - выхоодим из цикла чтения
		if readedURL == nil {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//	записываем список считанных URL в журнал
		log.Println(readedURL)
		//	добавляем связку HASH<==>URL в таблицу в RAM
		URLTable[readedURL.HashURL] = readedURL.LongURL
	}
}

// Метод для сохранения в БД связки короткого и длинного URL.
func Insert(shortURL, longURL string) {
	//	Проверяем наличие <shorten_URL> в списке сохраненных URL
	//	если такой URL уже есть в базе, то повторную вставку не производим
	_, ok := URLTable[shortURL]
	if !ok {
		URLTable[shortURL] = longURL
		//	если файл для хранения URL не задан, то храним список только в RAM
		if FilePath != "" {
			//	создаем экземпляр структуры хранения связки HASH<==>URL
			shortenURL := ShortenURL{
				HashURL: shortURL,
				LongURL: longURL,
			}
			//	создаем экземпляр WRITER в файл-хранилище HASH<==>URL
			writtenURL, err := NewURLWriter(FilePath)
			if err != nil {
				log.Fatal(err)
			}
			defer writtenURL.Close()
			//	производим запись в файл-хранилище связки HASH<==>URL
			if err := writtenURL.WriteURL(&shortenURL); err != nil {
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
