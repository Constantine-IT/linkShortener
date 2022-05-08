package main

import (
	"github.com/Constantine-IT/linkShortener/cmd/shortener/internal/handlers"
	"github.com/Constantine-IT/linkShortener/cmd/shortener/internal/storage"
	"net/http"
)

func main() {
	//	конфигурация приложения через считывание флагов и переменных окружения
	cfg := newConfig()

	//	инициализируем источники данных нашего приложения для работы с URL
	datasource, err := storage.NewDatasource(cfg.DatabaseDSN, cfg.FileStorage)
	if err != nil {
		cfg.ErrorLog.Fatal(err)
	}

	//	инициализируем контекст нашего приложения
	app := &handlers.Application{
		ErrorLog:   cfg.ErrorLog, //	журнал ошибок
		InfoLog:    cfg.InfoLog,  //	журнал информационных сообщений
		BaseURL:    cfg.BaseURL,  //	базоовый адрес сервера
		Datasource: datasource,   //	источник данных для хранения URL
	}

	//	при остановке сервера отложенно закроем все источники данных
	defer app.Datasource.Close()

	//	запуск сервера
	cfg.InfoLog.Printf("Server starts at address: %s", cfg.ServerAddress)

	srv := &http.Server{
		Addr:     cfg.ServerAddress,
		ErrorLog: cfg.ErrorLog,
		Handler:  app.Routes(),
	}
	cfg.ErrorLog.Fatal(srv.ListenAndServe())
}
