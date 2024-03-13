# Старт с базового образа golang
FROM golang:1.22.1-alpine AS builder

# Установка рабочей директории внутри контейнера
WORKDIR /src

# Копирование go mod и sum файлов
COPY ./ms-users/go.mod .
COPY ./ms-users/go.sum .

# Загрузка всех зависимостей.
# Если зависимости не изменились, то кэш слоя docker будет использован
RUN go mod download

# Копирование исходного кода из текущей директории в рабочую директорию внутри контейнера
COPY ./ms-users/ .

#  Copy env files
COPY ./ms-users/.env .env

# Компиляция приложения
RUN go build -o main .

# Вторая стадия сборки приложения
FROM alpine

# Установка рабочей директории внутри контейнера
WORKDIR /app

#  Копирование исходного кода из текущей стадии в рабочую директорию внутри контейнера
COPY --from=builder /src/main /app/

#  Copying env files
COPY --from=builder /src/.env .env

# Запуск бинарного файла
CMD ["./main"]