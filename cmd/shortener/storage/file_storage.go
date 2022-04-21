package storage

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"
)

//	Структуры и методы работы с файловым хранилищем URL

//	FileStorage - путь к файлу-хранилищу URL

var URLwriter *writer
var URLreader *reader

//	структура файлового дескриптора для записи
type writer struct {
	mutex   sync.Mutex
	file    *os.File
	encoder *json.Encoder
}

//	конструктор, создающий экземпляр файлового дескриптора для записи
func NewWriter(fileName string) (*writer, error) {
	//	файл открывается только на запись с добавлением в конец файла, если файла нет - создаем
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// метод записи в файл для экземпляра файлового дескриптора для записи
func (p *writer) Write(URL *shortenURL) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.encoder.Encode(&URL)
}

// метод закрытия файла для экземпляра файлового дескриптора для записи
func (p *writer) Close() error {
	return p.file.Close()
}

//	структура файлового дескриптора для чтения
type reader struct {
	mutex   sync.Mutex
	file    *os.File
	decoder *json.Decoder
}

//	конструктор, создающий экземпляр файлового дескриптора для чтения
func NewReader(fileName string) (*reader, error) {
	//	файл открывается только на чтение, если файла нет - создаем
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// метод чтения из файла для экземпляра файлового дескриптора для чтения
func (c *reader) Read() (*shortenURL, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	shortenURL := &shortenURL{}
	if err := c.decoder.Decode(&shortenURL); err != nil {
		return nil, err
	}
	return shortenURL, nil
}

// метод закрытия файла для экземпляра файлового дескриптора для чтения
func (c *reader) Close() error {
	return c.file.Close()
}

//	Метод первичного заполнения хранилища URL из файла сохраненных URL, при старте сервера
func InitialURLFulfilment(s *Storage) {

	//	блокируем хранилище URL в оперативной памяти на время заливки данных
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for {
		//	считываем записи по одной из файла-хранилища HASH + <original_URL> + UserID
		readURL, err := URLreader.Read()
		//	когда дойдем до конца файла - выходим из цикла чтения
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//	добавляем связку HASH и (<original_URL> + UserID) в хранилище
		s.Data[readURL.HashURL] = RowStorage{readURL.LongURL, readURL.UserID}
	}
}
