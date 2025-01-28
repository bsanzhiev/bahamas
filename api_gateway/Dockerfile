# Старт с базового образа golang
FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY ./ms-gateway/go.mod .
COPY ./ms-gateway/go.sum .
RUN go mod download
COPY ./ms-gateway/ .
RUN go build -o main .
FROM alpine
WORKDIR /app
COPY --from=builder /src/main /app/
COPY ./ms-gateway/ .
CMD ["./main"]
