package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	ServerAddress string
	BaseURL       string
	DatabaseDSN   string
	FileStorage   string
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
}

//	newConfig - функция-конфигуратор приложения через считывание флагов и переменных окружения
func newConfig() (cfg Config) {
	//	Приоритеты настроек:
	//	1.	Переменные окружения - ENV
	//	2.	Значения, задаваемые флагами при запуске из консоли
	//	3.	Значения по умолчанию.

	//	Считываем флаги запуска из командной строки и задаём значения по умолчанию, если флаг при запуске не указан
	ServerAddress := flag.String("a", "127.0.0.1:8080", "SERVER_ADDRESS - адрес запуска HTTP-сервера")
	BaseURL := flag.String("b", "http://127.0.0.1:8080", "BASE_URL - базовый адрес результирующего сокращённого URL")
	DatabaseDSN := flag.String("d", "", "DATABASE_DSN - адрес подключения к БД (PostgreSQL)")
	FileStorage := flag.String("f", "", "FILE_STORAGE_PATH - путь до файла с сокращёнными URL")
	//	парсим флаги
	flag.Parse()

	//	считываем переменные окружения
	//	если они заданы - переопределяем соответствующие локальные переменные:
	if u, flg := os.LookupEnv("SERVER_ADDRESS"); flg {
		*ServerAddress = u
	}
	if u, flg := os.LookupEnv("BASE_URL"); flg {
		*BaseURL = u
	}
	if u, flg := os.LookupEnv("DATABASE_DSN"); flg {
		*DatabaseDSN = u
	}
	if u, flg := os.LookupEnv("FILE_STORAGE_PATH"); flg {
		*FileStorage = u
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  // logger для информационных сообщений
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // logger для сообщений об ошибках

	// сигнальный канал для отслеживания системных вызовов на остановку агента
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	//	запускаем процесс слежение за сигналами на останов агента
	go func() {
		for {
			s := <-signalChanel
			if s == syscall.SIGINT || s == syscall.SIGTERM || s == syscall.SIGQUIT {
				cfg.InfoLog.Println("SERVER link shortener normal SHUTDOWN (code 0)")
				os.Exit(0)
			}
		}
	}()

	//	собираем конфигурацию агента
	cfg = Config{
		ServerAddress: *ServerAddress,
		BaseURL:       *BaseURL,
		DatabaseDSN:   *DatabaseDSN,
		FileStorage:   *FileStorage,
		InfoLog:       infoLog,
		ErrorLog:      errorLog,
	}

	//	выводим в лог конфигурацию сервера
	log.Println("SERVER link shortener STARTED with configuration:\n   SERVER_ADDRESS: ", cfg.ServerAddress, "\n   BASE_URL: ", cfg.BaseURL, "\n   DATABASE_DSN: ", cfg.DatabaseDSN, "\n   FILE_STORAGE_PATH: ", cfg.FileStorage)

	return cfg
}
