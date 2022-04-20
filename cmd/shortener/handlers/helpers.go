package handlers

import (
	"crypto/md5"
	"fmt"
	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
	"strings"
)

//	saveURLtoDB - вспомогательная функция, создающая HASH из связки (<original_URL> + UserID),
//	сохраняет строку (HASH + <original_URL> + UserID) и возвращает короткий URL для отправки клиенту
func (app *Application) saveURLtoDB(longURL, userID string) (string, error) {
	var err error
	// изготавливаем HASH из входящего URL с помощью MD5 hash algorithm
	mdSum := md5.Sum([]byte(longURL + userID))
	hash := fmt.Sprintf("%X", mdSum[0:4])

	//	если указан connect к базе данных, то работаем только в ней
	if app.Database != nil {
		//	проверяем, есть ли URL уже в нашей БД,
		//	если есть, то возвращаем уже сохраненный ранее <shorten_URL>, с ошибкой ErrConflictRecord
		preHash, flg := app.Database.GetByLongURL(longURL)
		if flg {
			app.InfoLog.Printf("%v - %s", storage.ErrConflictRecord, longURL)
			// Изготавливаем  <shorten_URL> из базового адреса нашего сервера и HASH
			shortURL := strings.Join([]string{app.BaseURL, preHash}, "/")
			return shortURL, storage.ErrConflictRecord
		}
		//	если такой записи в нашей БД нет, то добавляем её в нашу базу
		err = app.Database.Insert(hash, longURL, userID)
		if err != nil {
			app.ErrorLog.Println("New URL INSERT - FAILED: " + err.Error())
			return "", err
		} else {
			app.InfoLog.Println("New URL INSERT - SUCCESS - " + longURL)
		}
	}

	//	если база данных не задействуется, то работаем со структурами в RAM и файлом-хранилищем
	if app.Storage != nil {
		err = app.Storage.Insert(hash, longURL, userID, app.FileStorage)
	}
	if err == nil {
		// Изготавливаем  <shorten_URL> из базового адреса нашего сервера и HASH
		shortURL := strings.Join([]string{app.BaseURL, hash}, "/")
		app.InfoLog.Println("New SHORT URL was created - " + shortURL)
		return shortURL, nil
	} else {
		return "", err
	}
}
