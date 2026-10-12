#  MiniPubSub — High-Performance In-Memory PubSub Broker

Легковесный in-memory брокер сообщений по модели PublishSubscribe (Fan-Out), написанный на Go. Сервис обеспечивает изоляцию медленных подписчиков с помощью буферизованных каналов и неблокирующей отправки, предоставляет HTTP REST API для публикации и Server-Sent Events (SSE) для стриминга данных клиентам в реальном времени.

---

##  Архитектура и структура проекта

Проект спроектирован по канонам Standard Go Project Layout с четким разделением слоев абстракции

```text
minipubsub
├── cmd
│    └── main.go              # Точка входа сборка зависимостей, перехват OS-сигналов и Graceful Shutdown
│
├── internal                    # Приватный код приложения (бизнес-логика и адаптеры)
│   ├── broker                  # Ядро брокера (DomainCore)
│   │   ├── broker.go            # Интерфейс брокера, сущность Broker и управление топиками
│   │   ├── broker_test.go       # Параллельные тесты ядра и проверка на гонки (-race)
│   │   ├── subscriber.go        # Модель подписчика, жизненный цикл очередей
│   │   ├── topic.go             # Логика конкретной темы синхронизация и fan-out рассылка
│   │   └── topic_test.go        # Table-driven тесты операций SubscribeUnsubscribePublish
│   │
│   ├── config                  # Конфигурация сервиса
│   │   └── config.go            # Чтение параметров запуска (порт, размеры буферов каналов)
│   │
│   ├── handler                 # Транспортный уровень (HTTP Transport)
│   │   ├── http.go              # HTTP-хэндлеры REST API и Server-Sent Events (SSE)
│   │   └── http_test.go         # Интеграционные unit-тесты хэндлеров через httptest
│   │
│   └── model                   # Доменные структуры и DTO
│       └── message.go           # Схемы сообщений и сериализации входных данных
│
├── Dockerfile                   # Оптимизированная multi-stage сборка контейнера
├── Makefile                     # Автоматизация типовых команд (run, test, race, lint)
├── go.mod                       # Описание Go-модуля и зависимостей
├── go.sum                       # Контрольные суммы используемых библиотек
└── README.md                    # Документация проекта