package models

import "errors"

// DataSource - Определяем тип который обертывает подключения к хранилищу данных
// type DataSource struct {}
// Пока работаем с мапой, висящей в RAM
// Потом перепишем всё на работу с БД
var urlTable = make(map[string]string) // таблица для хранения URL

// Методы работы с моделью данных
// Insert - Метод для вставки в базу связки короткого и длинного URL.
func Insert(shortURL, longURL string) error {
	urlTable[shortURL] = longURL
	return nil
}

// Get - Метод для нахождения длинного URL по короткому URL.
func Get(shortURL string) (string, error) {

	longURL := urlTable[shortURL]
	if longURL == "" {
		return longURL, errors.New("URLs: подходящей записи не найдено")
	}
	return longURL, nil
}
