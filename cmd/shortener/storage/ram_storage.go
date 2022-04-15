package storage

import (
	"errors"
	"io"
	"log"
	"sync"
)

//	Структура хранилища URL в оперативной памяти
type rowStorage struct {
	longURL string
	userID  string
}

type Storage struct {
	data  map[string]rowStorage
	mutex sync.Mutex
}

//	Констуктор хранилища URL в оперативной памяти
func NewStorage() *Storage {
	return &Storage{data: make(map[string]rowStorage)}
}

//	Метод первичного заполнения хранилища URL из файла сохраненных URL, при старте сервера
func InitialFulfilmentURLDB(storage *Storage, file string) {
	//	создаем экземпляр reader для файла-хранилища HASH<==>URL
	fileReader, err := newReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fileReader.close()

	//	блокируем ресурс на время работы с ним
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	for {
		//	считываем записи по одной из файла-хранилища HASH<==>URL + UserID
		readURL, err := fileReader.read()
		//	когда дойдем до конца файла - выходим из цикла чтения
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//	добавляем связку HASH<==>URL + UserID в таблицу в RAM
		storage.data[readURL.HashURL] = rowStorage{readURL.LongURL, readURL.UserID}
	}
}
