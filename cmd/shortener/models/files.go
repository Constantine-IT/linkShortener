package models

import (
	"encoding/json"
	"os"
)

//	структура записи для сохраниния связки HASH<==>URL
type ShortenURL struct {
	HashURL string `json:"hash-url"`
	LongURL string `json:"long-url"`
}

//	структура файлового дескриптора для записи
type URLWriter struct {
	file    *os.File
	encoder *json.Encoder
}

//	конструктор, создающий экземпляр файлового дескриптора для записи
func NewURLWriter(fileName string) (*URLWriter, error) {
	//	файл открывается только на запись с добавлением в конец файла, если файла нет - создаем
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &URLWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// метод записи в файл для экземпляра файлового дескриптора для записи
func (p *URLWriter) WriteURL(URL *ShortenURL) error {
	return p.encoder.Encode(&URL)
}

// метод закрытия файла для экземпляра файлового дескриптора для записи
func (p *URLWriter) Close() error {
	return p.file.Close()
}

//	структура файлового дескриптора для чтения
type URLReader struct {
	file    *os.File
	decoder *json.Decoder
}

//	конструктор, создающий экземпляр файлового дескриптора для чтения
func NewURLReader(fileName string) (*URLReader, error) {
	//	файл открывается только на чтение, если файла нет - создаем
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &URLReader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// метод чтения из файла для экземпляра файлового дескриптора для чтения
func (c *URLReader) ReadURL() (*ShortenURL, error) {
	shortenURL := &ShortenURL{}
	if err := c.decoder.Decode(&shortenURL); err != nil {
		return nil, err
	}
	return shortenURL, nil
}

// метод закрытия файла для экземпляра файлового дескриптора для чтения
func (c *URLReader) Close() error {
	return c.file.Close()
}
