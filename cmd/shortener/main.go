package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	h "github.com/Constantine-IT/linkShortener/cmd/shortener/handlers"
	m "github.com/Constantine-IT/linkShortener/cmd/shortener/models"
)

func main() {

	//	Приоритеты настроек:
	//	1.	Переменные окружения - ENV
	//	2.	Значения, задаваемые флагами при запуске из консоли
	//	3.	Значения по умолчанию.

	//	Считываем флаги запуска из командной строки
	serverAddress := flag.String("a", "127.0.0.1:8080", "SERVER_ADDRESS - адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "http://127.0.0.1:8080", "BASE_URL - базовый адрес результирующего сокращённого URL")
	fileStoragePath := flag.String("f", "", "FILE_STORAGE_PATH - путь до файла с сокращёнными URL")
	flag.Parse()

	/*	На будущее - есть возможность использовать файл конфигурации.
		//	чтение файла конфигурации сервера
			config, err := os.ReadFile("server.cfg")
			if err == nil {
				log.Printf("Читаем файл server.cfg \n %s", config)
			} else {
				log.Println(err.Error())
			}
			//	парсинг считанной конфигурации
			//	ServerAddress в server.cfg - адрес для запуска сервера
			_, err = fmt.Sscanf(string(config), "ServerAddress %s", &srvAddr)
			if err != nil {
				log.Println(err.Error())
			}
	*/
	srvAddr := *serverAddress
	h.Addr = *baseURL
	m.FilePath = *fileStoragePath

	//	считываем переменные окружения:
	if u, flag := os.LookupEnv("SERVER_ADDRESS"); flag {
		srvAddr = u
	}

	if u, flag := os.LookupEnv("BASE_URL"); flag {
		h.Addr = u
	}

	//	Путь к файлу для сохранения URL - FILE_STORAGE_PATH (при перезапуске сервера данные сохраняются)
	//	если FILE_STORAGE_PATH не задана, то храним URL только в оперативной памяти и теряем при перезапуске.
	if u, flag := os.LookupEnv("FILE_STORAGE_PATH"); flag {
		m.FilePath = u
	}
	//	m.FilePath = "url_DB.txt"	//	Это для тестов с файлом для хранения URL. На бою закомментировать.
	//	считываем файл сохраненных URL и заполняем этой информацией БД <shorten_URL>
	if m.FilePath != "" {
		m.InitialFulfilmentURLDB()
	}

	log.Printf("Сервер будет запущен по адресу: %s", srvAddr)

	//	запуск сервера
	srv := &http.Server{
		Addr:    srvAddr,
		Handler: h.Routes(),
	}
	log.Fatal(srv.ListenAndServe())
}
