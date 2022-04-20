package handlers

import "net/http"

//	PingDataBaseHandler - обработчик проверки доступности базы данных через GET /ping
func (app *Application) PingDataBaseHandler(w http.ResponseWriter, r *http.Request) {
	if app.Database == nil {
		http.Error(w, "DataBase environment variable wasn't set", http.StatusInternalServerError)
		app.ErrorLog.Println("Attempt to PING database, that wasn't set in server configuration")
		return
	}
	err := app.Database.DB.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		app.ErrorLog.Println("Try to PING database error: " + err.Error())
	} else {
		http.Error(w, http.StatusText(200), http.StatusOK)
	}
}
