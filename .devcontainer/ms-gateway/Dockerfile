# Старт с базового образа golang
FROM golang:1.22.1-alpine AS builder
# Установка рабочей директории внутри контейнера
WORKDIR /src
# Копирование go mod и sum файлов
COPY ../../ms-gateway/go.mod .
COPY ../../ms-gateway/go.sum .
# Загрузка всех зависимостей.
# Если зависимости не изменились, то кэш слоя docker будет использован
RUN go mod download
# Копирование исходного кода из текущей директории в рабочую директорию внутри контейнера
COPY ../../ms-gateway/ .
# Компиляция приложения
RUN go build -o main .
#
# Вторая стадия сборки приложения
FROM alpine
# Установка рабочей директории внутри контейнера
WORKDIR /app
#  Копирование бинароного файла
COPY --from=builder /src/main /app/build/

# Копирование исходного кода
# ERROR: failed to solve: failed to compute cache key: failed to calculate checksum of ref 
# d1d1fef4-6c6d-436a-a0b5-c6bf6c08f633::vcxpp3v09hpritqyvmq1uhmxq: "/ms-gateway": not found
COPY ../../ms-gateway/ .

# Запуск бинарного файла
CMD ["./main"]
