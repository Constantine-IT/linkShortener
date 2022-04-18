package handlers

import (
	"crypto/md5"
	"fmt"
	"github.com/Constantine-IT/linkShortener/cmd/shortener/storage"
	"strings"
)

//	saveURLtoDB - вспомогательная функция, создающая HASH из связки (URL + UserID),
//	сохраняет связку HASH<==>URL+UserID и возвращает короткий URL для отправки клиенту
func (app *Application) saveURLtoDB(longURL, userID string) (string, error) {
	var err error
	// изготавливаем HASH из входящего URL с помощью MD5 hash algorithm
	mdSum := md5.Sum([]byte(longURL + userID))
	hash := fmt.Sprintf("%X", mdSum[0:4])

	//	вызов метода-вставки в структуру хранения связки HASH<==>URL+UserID
	//	если указан connect к базе данных, то работаем только в ней
	if app.Database != nil {
		//	проверяем, есть ли URL уже в нашей БД,
		//	и если есть, то возвращаем уже сохраненный ранее <shorten_URL>, с ошибкой ErrConflictRecord
		prevHash, flg := app.Database.GetByLongURL(longURL)
		if flg {
			app.ErrorLog.Printf("%w - %s", storage.ErrConflictRecord, longURL)
			// Изготавливаем  <shorten_URL> из базового адреса нашего сервера и HASH
			shortURL := strings.Join([]string{app.BaseURL, prevHash}, "/")
			return shortURL, storage.ErrConflictRecord
		}
		//	если такой записи в нашей БД нет, то добавляем её в нашу базу
		err = app.Database.Insert(hash, longURL, userID)
	}

	//	если база данных не задействуется, то работаем со структурами в RAM и файлом-хранилищем
	if app.Storage != nil {
		err = app.Storage.Insert(hash, longURL, userID, app.FileStorage)
	}
	if err == nil {
		// Изготавливаем  <shorten_URL> из базового адреса нашего сервера и HASH
		shortURL := strings.Join([]string{app.BaseURL, hash}, "/")
		return shortURL, nil
	} else {
		return "", err
	}
}
