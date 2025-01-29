# Architecture

```
customers/
├── cmd/
│   └── main.go                 # Точка входа (инициализация конфига, запуск серверов)
├── internal/
│   ├── domain/                 # Ядро: бизнес-логика и сущности
│   │   ├── customer.go         # Сущность Customer (поля + бизнес-методы)
│   │   ├── events.go           # Доменные события (CustomerCreated, CustomerUpdated)
│   │   └── repository.go       # Интерфейсы: CustomerRepository (все методы БД)
│   │
│   ├── application/            # Сценарии использования и orchestration
│   │   ├── usecases/
│   │   │   ├── create.go       # CreateCustomerUseCase (валидация, вызов репозитория)
│   │   │   ├── update.go       # UpdateCustomerUseCase
│   │   │   └── saga_handlers/  # Обработчики событий саг
│   │   │       └── kyc_saga.go # Компенсация при откате KYC
│   │   │
│   │   └── events/             # Обработчики доменных событий
│   │       └── created.go      # Публикация CustomerCreated в Kafka
│   │
│   ├── infrastructure/         # Реализация внешних зависимостей
│   │   ├── postgres/           # Реализация CustomerRepository для PostgreSQL
│   │   │   └── repository.go
│   │   ├── kafka/              # Producer для публикации событий
│   │   │   └── publisher.go
│   │   └── grpc/               # gRPC-сервер (без бизнес-логики)
│   │       ├── server.go
│   │       └── converter.go    # Конвертация gRPC <-> Domain
│   │
│   └── interfaces/             # Адаптеры для внешних систем
│       ├── grpc_handler/       # Обработчики gRPC-запросов
│       │   └── customer.go     # (Вызывает Use Cases)
│       └── message_consumer/   # Потребители событий из брокеров
│           ├── kafka_consumer.go
│           └── handlers/
│               └── loan_approved.go # Обработчик события "LoanApproved"
├── configs/
│   ├── config.yaml             # Конфиг (порты, URL БД, топики Kafka)
│   └── env.go                  # Загрузка конфига из переменных окружения
├── migrations/                 # Миграции PostgreSQL
│   └── 0001_init_customers.sql
├── proto/                      # Сгенерированные gRPC-стабы (из общей папки libs/protobuf)
├── go.mod                      # Зависимости модуля
└── Dockerfile                  # Мультистейдж-сборка

```