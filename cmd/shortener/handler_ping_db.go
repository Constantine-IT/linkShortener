package main

import "net/http"

//	PingDataBaseHandler - обработчик проверки доступности базы данных через GET /ping
func (app *application) PingDataBaseHandler(w http.ResponseWriter, r *http.Request) {
	if app.database == nil {
		http.Error(w, "DataBase environment variable wasn't set", http.StatusInternalServerError)
		return
	}
	err := app.database.DB.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		http.Error(w, http.StatusText(200), http.StatusOK)
	}
}
