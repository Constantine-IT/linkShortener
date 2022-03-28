package models

// Пока работаем с MAP, висящей в RAM, потом перепишем всё на работу с БД
var urlTable = make(map[string]string) // таблица для хранения URL

// Методы работы с моделью данных

// Insert - Метод для вставки в базу связки короткого и длинного URL.
func Insert(shortURL, longURL string) {
	urlTable[shortURL] = longURL
}

// Get - Метод для нахождения длинного URL по короткому URL.
func Get(shortURL string) (longURL string, flag bool) {
	longURL, ok := urlTable[shortURL]
	return longURL, ok
}
