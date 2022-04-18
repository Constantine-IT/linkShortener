package main

import (
	"crypto/md5"
	"fmt"
	"strings"
)

//	saveURLtoDB - вспомогательная функция, создающая HASH из связки (URL + UserID),
//	сохраняет связку HASH<==>URL+UserID и возвращает короткий URL для отправки клиенту
func (app *application) saveURLtoDB(longURL, userID string) (string, error) {
	var err error
	// изготавливаем HASH из входящего URL с помощью MD5 hash algorithm
	mdSum := md5.Sum([]byte(longURL + userID))
	hash := fmt.Sprintf("%X", mdSum[0:4])

	// вызов метода-вставки в структуру хранения связки HASH<==>URL+UserID
	if app.storage != nil {
		err = app.storage.Insert(hash, longURL, userID, app.fileStorage)
	}
	if app.database != nil {
		err = app.database.Insert(hash, longURL, userID)
	}

	if err == nil {
		// Изготавливаем  <shorten_URL> из базового адреса нашего сервера и HASH
		shortURL := strings.Join([]string{app.baseURL, hash}, "/")
		return shortURL, nil
	} else {
		return "", err
	}
}
