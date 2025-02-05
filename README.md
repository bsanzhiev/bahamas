# Architecture

## Project structure

```
online-bank-monorepo/
├── .github/                   # GitHub Actions workflows
│   └── workflows/
│       ├── ci.yaml            # CI для тестов и линтеров
│       └── cd.yaml            # CD для деплоя в Kubernetes
├── api_gateway/               # REST → gRPC шлюз
│   ├── cmd/
│   │   └── main.go            # Точка входа
│   ├── internal/
│   │   ├── handlers/          # REST-обработчики
│   │   ├── grpc_clients/      # Клиенты для микросервисов
│   │   └── middleware/        # Аутентификация, логирование
│   ├── configs/
│   │   ├── routes.yaml        # Маршрутизация REST → gRPC
│   │   └── config.yaml        # Настройки окружения
│   ├── proto/                 # Сгенерированные gRPC-стабы (из общей папки)
│   ├── Dockerfile             # Мультистейдж-сборка
│   └── go.mod                 # Зависимости
├── services/                  # Все микросервисы
│   ├── customers/
│   │   ├── cmd/
│   │   │   └── main.go        # Запуск сервиса
│   │   ├── internal/
│   │   │   ├── repository/    # Репозиторий PostgreSQL
│   │   │   ├── service/       # Бизнес-логика
│   │   │   └── grpc/          # gRPC-сервер
│   │   ├── migrations/        # SQL-миграции (goose)
│   │   │   └── 001_init.sql
│   │   ├── configs/
│   │   │   └── config.yaml    # Настройки БД и портов
│   │   ├── proto/             # Сгенерированные .pb.go файлы
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── transactions/          # Аналогичная структура
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── migrations/        # Cassandra миграции (если нужно)
│   │   └── ...
│   └── ...                    # accounts, cards, loans и др.
├── libs/                      # Общие библиотеки
│   ├── protobuf/              # Единые .proto-файлы
│   │   ├── transactions/
│   │   │   └── transaction.proto
│   │   ├── customers/
│   │   │   └── customer.proto
│   │   └── ...
│   ├── auth/                  # Общая аутентификация
│   │   ├── jwt.go
│   │   └── middleware.go
│   ├── kafka_client/          # Обертка для Sarama
│   │   └── producer.go
│   └── ...                    # Другие утилиты
├── infrastructure/
│   ├── kubernetes/            # Манифесты для k8s
│   │   ├── api-gateway/
│   │   │   ├── deployment.yaml
│   │   │   └── service.yaml
│   │   ├── customers/
│   │   │   ├── deployment.yaml
│   │   │   └── pvc.yaml       # Persistent Volume для PostgreSQL
│   │   └── ...
│   ├── docker-compose/        # Локальное окружение
│   │   ├── docker-compose.yaml
│   │   └── .env.example       # Переменные окружения
│   └── terraform/             # Облачная инфраструктура
│       ├── modules/
│       │   ├── network/       # VPC, подсети
│       │   └── eks/           # EKS кластер
│       └── environments/
│           ├── prod/
│           └── staging/
├── message_brokers/           # Конфиги брокеров
│   ├── kafka/
│   │   ├── docker-compose.yaml
│   │   └── topics.sh          # Скрипт создания топиков
│   └── rabbitmq/
│       └── definitions.json   # Очереди и пользователи
├── saga_orchestrator/         # Оркестратор распределенных транзакций
│   ├── cmd/
│   │   └── main.go            # Запуск оркестратора
│   ├── internal/
│   │   ├── workflows/         # Сценарии саг (например, кредит)
│   │   └── compensation/      # Компенсирующие действия
│   └── go.mod
├── scripts/                   # Вспомогательные скрипты
│   ├── generate_proto.sh      # Генерация gRPC-кода из .proto
│   ├── migrate_db.sh          # Запуск миграций для сервиса
│   └── deploy/                # Утилиты для деплоя
├── docs/
│   ├── API.md                 # OpenAPI спецификация для REST
│   ├── ARCHITECTURE.md        # Диаграммы C4 и последовательностей
│   └── DEVELOPMENT.md         # Инструкции для разработчиков
├── .gitignore
├── go.work                    # Go Workspace (для монорепы)
├── Makefile                   # Основные команды
└── README.md                  # Общее описание проекта
```

## Customers Service structure

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

get customer: 
1. Клиент отправляет HTTP-запрос (например, через REST API) в api_gateway.
2. В api_gateway/cmd/main.go инициализируется HTTP-обработчик (customerHandler), который отвечает за обработку входящих HTTP-запросов.
```go
	customerHandler := http_handler.NewCustomerHandler(*customerService, logger)
```
3. HTTP-обработчик (customerHandler) принимает HTTP-запрос, парсит его и вызывает соответствующий метод в customerService.
4. customerService — это слой бизнес-логики (Use Case), который знает, как обработать запрос и взаимодействовать с микросервисом customers.
5. customerService преобразует данные из HTTP-запроса в формат, понятный для gRPC.
6. Затем он отправляет gRPC-запрос в микросервис customers.
7. В микросервисе customers есть gRPC-обработчик (grpc_handler/customer.go), который принимает запрос.
8. Этот обработчик вызывает соответствующий метод в слое бизнес-логики (Use Case) микросервиса customers.
9. Use Case в микросервисе customers выполняет необходимую бизнес-логику (например, получение данных о клиенте из базы данных).
10. Результат возвращается обратно в gRPC-обработчик.
11. Микросервис customers отправляет gRPC-ответ обратно в api_gateway.
12. customerService в api_gateway получает ответ и преобразует его в HTTP-формат.
13. HTTP-обработчик (customerHandler) формирует HTTP-ответ и отправляет его клиенту.
