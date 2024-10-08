# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Создаем моки для тестов слоев сервиса
mockgen -destination=internal/server/controller/http/mock_service.go -package=http  github.com/imirjar/metrx/internal/server/controller/http Service

mockgen -destination=internal/server/service/mock_storage.go -package=service  github.com/imirjar/metrx/internal/server/service Storager

mockgen -destination=internal/agent/app/mock_client.go -package=app  github.com/imirjar/metrx/internal/agent/app Client

## Проверить процент покрытия можно с помощью команды:
go test -coverprofile=coverage.out ./... ;    go tool cover -func=coverage.out

На данный момент он составляет:
total:         (statements)            18.8%


## Контейнер с тестовой базой
sudo docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=praktikum -d postgres:latest



## RSA
используем .pem формат для хранения ключей
ключи обязятельно должны быть достаточной длинны для шифрования большого количества отправляемых метрик
я использую размер 4096


## PROTOC
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  internal/api/api.proto 