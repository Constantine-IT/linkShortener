# cmd/shortener
В данной директории будет содержаться код, который скомпилируется в бинарное приложение.

Конфигурация для запуска сервера в файле server.cfg

Запускаем сервер командой go run .

# go-musthave-shortener-tpl
Шаблон репозитория для практического трека «Go в веб-разработке».

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` - адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

# PUSH в репозиторий GitHub:

Сначала создаем новую ветку:

git checkout -b increment<n>

Помечаем всё что менялось с последнего PUSH:

git add .

git commit -m 'increment<n>'

Ну и делаем PUSH в Github:

git push --set-upstream origin increment<n>
  
# В проекте используются доп. библиотеки TESTIFY и CHI

Их надо предварительно зарегистрировать в модуле:

go get github.com/go-chi/chi/v5  

go get github.com/go-chi/chi/v5/middleware

go get github.com/stretchr/testify/assert

go get github.com/stretchr/testify/require
  
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
