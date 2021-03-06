package handlers

import (
	"net/http"

	"github.com/Constantine-IT/linkShortener/cmd/shortener/internal/storage"
)

//	PingDataBaseHandler - обработчик проверки доступности базы данных через GET /ping
func (app *Application) PingDataBaseHandler(w http.ResponseWriter, r *http.Request) {
	switch value := app.Datasource.(type) {
	//	проверяем тип источника данных - Datasource
	case *storage.Database: //	если интерфейс источника данных Datasource реализован базой данных - Database
		err := value.DB.Ping() //	проверяем доступность базы данных
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			app.ErrorLog.Println("Try to PING database error: " + err.Error())
		} else {
			http.Error(w, http.StatusText(200), http.StatusOK)
		}
	default: //	если интерфейс источника данных Datasource реализован иной структурой
		http.Error(w, "DataBase environment variable wasn't set", http.StatusInternalServerError)
		app.ErrorLog.Println("Attempt to PING database, that wasn't set in server configuration")
		return
	}

}
