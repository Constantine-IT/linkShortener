package storage

import (
	"errors"
	"io"
	"log"
	"sync"
)

//	Структура хранилища URL в оперативной памяти
type Storage struct {
	data map[string]string
	mu   sync.Mutex
}

//	Констуктор хранилища URL в оперативной памяти
func NewStorage() *Storage {
	return &Storage{data: make(map[string]string)}
}

//	Метод первичного заполнения хранилища URL из файла сохраненных URL, при старте сервера
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
