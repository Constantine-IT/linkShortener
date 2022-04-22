package handlers

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
)

//	saveURLtoDB - вспомогательная функция, создающая HASH из связки (<original_URL> + UserID),
//	сохраняет строку (HASH + <original_URL> + UserID) и возвращает короткий URL для отправки клиенту
func (app *Application) saveURLtoDB(longURL, userID string) (string, error) {
	var err error
	// изготавливаем HASH из входящего URL с помощью MD5 hash algorithm
	mdSum := md5.Sum([]byte(longURL + userID))
	hash := fmt.Sprintf("%X", mdSum[0:4])

	//	проверяем, есть ли URL уже в нашей базе данных
	//	если есть, то возвращаем сохраненный ранее <shorten_URL>, с ошибкой ErrConflictRecord
	preHash, flg := app.Datasource.GetByLongURL(longURL)
	if flg {
		app.InfoLog.Printf("%v - %s", storage.ErrConflictRecord, longURL)
		// Изготавливаем  <shorten_URL> из базового адреса нашего сервера и HASH
		shortURL := strings.Join([]string{app.BaseURL, preHash}, "/")
		return shortURL, storage.ErrConflictRecord
	}
	//	если такой записи в нашей БД нет, то добавляем её в нашу базу
	err = app.Datasource.Insert(hash, longURL, userID)

	if err == nil {
		// Изготавливаем  <shorten_URL> из базового адреса нашего сервера и HASH
		shortURL := strings.Join([]string{app.BaseURL, hash}, "/")
		return shortURL, nil
	} else {
		return "", err
	}
}
