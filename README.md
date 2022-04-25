# Сервис сокращения ссылок - LinkShortener
# Конфигурация сервера через переменные окружения:
1. SERVER_ADDRESS - адрес запуска HTTP-сервера (по умолчанию - 127.0.0.1:8080)
2. BASE_URL - базовый адрес сокращённого URL (по умолчанию - http://127.0.0.1:8080)
3. DATABASE_DSN - адрес подключения к БД (PostgreSQL v.10+) (по умолчанию - "", т.е. работаем без БД)
3. FILE_STORAGE_PATH - файл хранения сокращённых URL (по умолчанию - "", т.е. работаем без файла)

`/cmd/shortener` - здесь содержится код, который компилируется в бинарное приложение.

Запускаем сервер командой `go run .`

доступны флаги запуска:
1. флаг -a, отвечающий за адрес запуска HTTP-сервера (переменная SERVER_ADDRESS);
2. флаг -b, отвечающий за базовый адрес сокращённого URL (переменная BASE_URL);
3. флаг -d, отвечающий за адрес подключения к БД (PostgreSQL v.10+) (переменная DATABASE_DSN)
4. флаг -f, отвечающий за путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH).

`Переменные среды имеют приоритет над флагами.`

# handlers
пакет содержит все маршруты web-сервера, их обработчики и помощники.
# storage
пакет содержит все структуры данных и методы для работы с ними, поддерживаются хранилища:
- в оперативной памяти (RAM),
- в текстовом файле,
- на сервере баз данных.

# В проекте используются доп. библиотеки TESTIFY и CHI
Их надо предварительно зарегистрировать в модуле:
```
go get github.com/go-chi/chi/v5  
go get github.com/go-chi/chi/v5/middleware
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
```
