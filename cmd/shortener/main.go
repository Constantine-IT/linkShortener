package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	h "github.com/Constantine-IT/linkShortener/cmd/shortener/handlers"
	m "github.com/Constantine-IT/linkShortener/cmd/shortener/models"
)

func main() {
	//	по умолчанию запускаем сервер на адресе 127.0.0.1:8080
	srvAddr := "127.0.0.1:8080"
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

	//	считываем переменные окружения:
	//	адрес запуска HTTP-сервера - SERVER_ADDRESS (default: 127.0.0.1:8080)
	if u, flag := os.LookupEnv("SERVER_ADDRESS"); flag {
		srvAddr = u
	}

	//	адрес для формирования <shorten_URL> - BASE_URL (default: http://127.0.0.1:8080)
	if u, flag := os.LookupEnv("BASE_URL"); flag {
		h.Addr = u
	}

	//	Путь к файлу для сохранения URL - FILE_STORAGE_PATH (при перезапуске сервера данные сохраняются)
	//	если FILE_STORAGE_PATH не задана, то храним URL только в оперативной памяти и теряем при перезапуске.
	if u, flag := os.LookupEnv("FILE_STORAGE_PATH"); flag {
		m.FilePath = u
	}
	//	m.FilePath = "url_DB.txt"	//	Это для тестов с файлом для хранения URL. На бою закомментировать.
	if m.FilePath != "" {
		fileReader, err := m.NewURLReader(m.FilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer fileReader.Close()
		log.Println("Из файла считаны сохраненные URL:")
		for {
			readedURL, err := fileReader.ReadURL()
			if readedURL == nil {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Println(readedURL)
			m.URLTable[readedURL.HashURL] = readedURL.LongURL
		}
	}

	log.Printf("Сервер будет запущен по адресу: %s", srvAddr)

	//	запуск сервера
	srv := &http.Server{
		Addr:    srvAddr,
		Handler: h.Routes(),
	}
	log.Fatal(srv.ListenAndServe())
}
