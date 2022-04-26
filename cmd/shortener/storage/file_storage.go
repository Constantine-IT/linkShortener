package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
)

//	Структуры и методы работы с файлом-хранилищем URL

//	fileWriter и fileReader - рабочие экземпляры файловых дескрипторов чтения и записи
var fileWriter *writer
var fileReader *reader

//	writer - структура файлового дескриптора для записи
type writer struct {
	mutex   sync.Mutex
	file    *os.File
	encoder *json.Encoder
}

//	NewWriter - конструктор, создающий экземпляр файлового дескриптора для записи
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

// Write - метод записи в файл для экземпляра файлового дескриптора для записи
func (p *writer) Write(URL *shortenURL) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.encoder.Encode(&URL)
}

// Close - метод закрытия файла для экземпляра файлового дескриптора для записи
func (p *writer) Close() error {
	return p.file.Close()
}

//	reader - структура файлового дескриптора для чтения
type reader struct {
	mutex   sync.Mutex
	file    *os.File
	decoder *json.Decoder
}

//	NewReader - конструктор, создающий экземпляр файлового дескриптора для чтения
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

// Read - метод чтения из файла для экземпляра файлового дескриптора для чтения
func (c *reader) Read() (*shortenURL, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	shortenURL := &shortenURL{}
	if err := c.decoder.Decode(&shortenURL); err != nil {
		return nil, err
	}
	return shortenURL, nil
}

// Close - метод закрытия файла для экземпляра файлового дескриптора для чтения
func (c *reader) Close() error {
	return c.file.Close()
}

//	InitialURLFulfilment - метод первичного заполнения хранилища URL из файла сохраненных URL, при старте сервера
func InitialURLFulfilment(s *Storage) error {

	//	блокируем хранилище URL в оперативной памяти на время заливки данных
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for {
		//	считываем записи по одной из файла-хранилища HASH + <original_URL> + UserID
		readURL, err := fileReader.Read()
		//	когда дойдем до конца файла - выходим из цикла чтения
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		//	добавляем связку HASH и (<original_URL> + UserID) в хранилище
		s.Data[readURL.HashURL] = RowStorage{readURL.LongURL, readURL.UserID}
	}
	return nil
}
