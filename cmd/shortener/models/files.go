package models

import (
	"encoding/json"
	"os"
)

type ShortenURL struct {
	HashURL string `json:"hash-url"`
	LongURL string `json:"long-url"`
}
type URLWriter struct {
	file    *os.File
	encoder *json.Encoder
}

func NewURLWriter(fileName string) (*URLWriter, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &URLWriter{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *URLWriter) WriteURL(URL *ShortenURL) error {
	return p.encoder.Encode(&URL)
}

func (p *URLWriter) Close() error {
	return p.file.Close()
}

type URLReader struct {
	file    *os.File
	decoder *json.Decoder
}

func NewURLReader(fileName string) (*URLReader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &URLReader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *URLReader) ReadURL() (*ShortenURL, error) {
	shortenURL := &ShortenURL{}
	if err := c.decoder.Decode(&shortenURL); err != nil {
		return nil, err
	}
	return shortenURL, nil
}

func (c *URLReader) Close() error {
	return c.file.Close()
}
