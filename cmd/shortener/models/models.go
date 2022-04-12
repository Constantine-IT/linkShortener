package models

import (
	"errors"
	"io"
	"log"
	"sync"
)

//	Структура хранилища URL
type Storage struct {
	data map[string]string
	mu   sync.Mutex
}

//	Констуктор хранилища URL
func NewStorage() *Storage {
	return &Storage{data: make(map[string]string)}
}

// Методы работы с хранилищем URL

//	Метод первичного заполнения БД из файла сохраненных URL при старте сервера
func InitialFulfilmentURLDB(storage *Storage, file string) {
	//	создаем экземпляр reader для файла-хранилища HASH<==>URL
	fileReader, err := newReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fileReader.close()
	log.Println("Обнаружен файл сохраненных URL. Начинаем считывание:")
	for {
		//	считываем записи по одной из файла-хранилища HASH<==>URL
		readURL, err := fileReader.read()
		//	когда дойдем до конца файла - выходим из цикла чтения
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//	записываем список считанных URL в log
		log.Println(readURL)
		//	добавляем связку HASH<==>URL в таблицу в RAM
		storage.mu.Lock()
		defer storage.mu.Unlock()
		storage.data[readURL.HashURL] = readURL.LongURL
	}
}

// Insert - Метод для сохранения в БД связки короткого и длинного URL.
func Insert(shortURL, longURL, filePath string, storage *Storage) error {
	//	пустые значения URL к вставке в хранилище не допускаются
	if shortURL == "" || longURL == "" {
		return errors.New("empty value is not allowed")
	}
	//	Проверяем наличие <shorten_URL> в списке сохраненных URL
	//	если такой URL уже есть в базе, то повторную вставку не производим
	if _, ok := storage.data[shortURL]; !ok {
		storage.mu.Lock()
		defer storage.mu.Unlock()
		storage.data[shortURL] = longURL
		//	если файл для хранения URL не задан, то храним список только в RAM в URLTable
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
