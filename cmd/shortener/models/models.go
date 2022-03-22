package models

import "errors"

// DataSource - Определяем тип который обертывает подключения к хранилищу данных
// type DataSource struct {}
// Пока работаем с мапой, висящей в RAM
// Потом перепишем всё на работу с БД
var urlTable = make(map[string]string) // таблица для хранения URL

// Методы работы с моделью данных
// Insert - Метод для создания связки короткого и длинного URL.
func Insert(shortUrl, longUrl string) error {
	urlTable[shortUrl] = longUrl
	return nil
}

// Get - Метод для возвращения длинного URL поего идентификатору ID.
func Get(shortUrl string) (string, error) {

	longUrl := urlTable[shortUrl]
	if longUrl == "" {
		return longUrl, errors.New("URLs: подходящей записи не найдено")
	}
	return longUrl, nil
}
