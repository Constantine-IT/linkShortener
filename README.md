# Сервис сокращения ссылок - LinkShortener
`/cmd/shortener`
В данной директории содержится код, который компилируется в бинарное приложение.
# Конфигурация сервера через переменные окружения:
1. SERVER_ADDRESS - адрес запуска HTTP-сервера (по умолчанию - 127.0.0.1:8080)
2. BASE_URL - базовый адрес сокращённого URL (по умолчанию - http://127.0.0.1:8080)
3. DATABASE_DSN - адрес подключения к БД (PostgreSQL v.10+) (по умолчанию - "", т.е. работаем без БД)
3. FILE_STORAGE_PATH - файл с сокращёнными URL (по умолчанию - "", т.е. работаем без файла, только в RAM)

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
пакет содержит все структуры данных и методы для работы с ними, с использованием разнообразных хранилищ
- в оперативной памяти (RAM),
- в файле-хранилище,
- на сервере баз данных.

# PUSH в репозиторий GitHub:
Сначала создаем новую ветку:
```
git checkout -b increment-N
```
Помечаем всё что менялось с последнего PUSH:
```
git add .
git commit -m 'increment-N'
```
Делаем PUSH в Github:
```
git push --set-upstream origin increment-N
```  
# В проекте используются доп. библиотеки TESTIFY и CHI
Их надо предварительно зарегистрировать в модуле:
```
go get github.com/go-chi/chi/v5  
go get github.com/go-chi/chi/v5/middleware
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
```
# Начало работы
1. `go-musthave-shortener-tpl` - шаблон репозитория для практического трека «Go в веб-разработке».
Создайте на его базе репозиторий в своём GitHub. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>`
(где `<name>` - адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.
# Обновление шаблона
Чтобы иметь возможность получать обновления автотестов и других частей шаблона выполните следующую команду:
```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-shortener-tpl.git
```
Для обновления кода автотестов выполните команду:
```
git fetch template && git checkout template/main .github
```
Затем добавьте полученные изменения в свой репозиторий.
