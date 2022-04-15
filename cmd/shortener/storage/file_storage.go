package storage

import (
	"encoding/json"
	"os"
	"sync"
)

//	Структуры и методы работы с файловым хранилищем URL

//	структура записи для сохраниния связки HASH<==>URL
type shortenURL struct {
	HashURL string `json:"hash-url"`
	LongURL string `json:"long-url"`
	UserID  string `json:"user-id"`
}

//	структура файлового дескриптора для записи
type writer struct {
	mutex   sync.Mutex
	file    *os.File
	encoder *json.Encoder
}

//	конструктор, создающий экземпляр файлового дескриптора для записи
func newWriter(fileName string) (*writer, error) {
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
func (p *writer) write(URL *shortenURL) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.encoder.Encode(&URL)
}

// метод закрытия файла для экземпляра файлового дескриптора для записи
func (p *writer) close() error {
	return p.file.Close()
}

//	структура файлового дескриптора для чтения
type reader struct {
	mutex   sync.Mutex
	file    *os.File
	decoder *json.Decoder
}

//	конструктор, создающий экземпляр файлового дескриптора для чтения
func newReader(fileName string) (*reader, error) {
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
func (c *reader) read() (*shortenURL, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	shortenURL := &shortenURL{}
	if err := c.decoder.Decode(&shortenURL); err != nil {
		return nil, err
	}
	return shortenURL, nil
}

// метод закрытия файла для экземпляра файлового дескриптора для чтения
func (c *reader) close() error {
	return c.file.Close()
}
