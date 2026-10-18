# MiniPubSub — High-Performance In-Memory Pub/Sub Broker

Легковесный in-memory брокер сообщений по модели **Publish/Subscribe (Fan-Out)**, написанный на Go. Сервис обеспечивает изоляцию медленных подписчиков с помощью буферизованных каналов и неблокирующей отправки (`select` + `default`), предоставляет HTTP REST API для публикации и Server-Sent Events (SSE) для стриминга данных клиентам в реальном времени.

---

## Стек технологий

* **Язык:** Go `1.24+` (горутины, буферизованные каналы, `sync.RWMutex`, `sync.Once`, `sync/atomic`, `context`)
* **Сетевой транспорт:** HTTP/1.1 (`net/http`), REST API, Server-Sent Events (`http.Flusher`)
* **Контейнеризация:** Docker (Multi-stage сборка на базе `golang:1.24-alpine` и `alpine:latest`), Docker Compose (лимит памяти `300M`, `GOMEMLIMIT=256MiB`)
* **Автоматизация:** GNU Make

---

## Архитектура и структура проекта

Проект организован по канонам Standard Go Project Layout с изоляцией слоев приложения:

```text
minipubsub/
├── cmd/
│   └── main.go                  # Точка входа, сборка зависимостей, Graceful Shutdown
│
├── internal/                    # Приватный код приложения (бизнес-логика и адаптеры)
│   ├── broker/                  # Ядро брокера (Domain/Core)
│   │   ├── broker.go            # Фасад брокера, реестр и управление топиками
│   │   ├── broker_test.go       # Параллельные стресс-тесты ядра, проверка на гонки (-race), бенчмарки
│   │   ├── subscriber.go        # Модель подписчика, жизненный цикл очередей
│   │   └── topic.go             # Логика темы, синхронизация и fan-out рассылка
│   │
│   ├── error/
│   │   └── error.go             # Ошибки уровня предметной области
│   │
│   ├── handler/                 # Транспортный уровень (HTTP Transport)
│   │   ├── http.go              # Обработчики REST API и Server-Sent Events (SSE)
│   │   └── http_test.go         # Интеграционные unit-тесты хэндлеров 
│   │
│   └── model/                   # Доменные структуры и DTO
│       └── message.go           # Схемы сообщений и сериализации входных данных
│
├── Dockerfile                   # Оптимизированная сборка контейнера (Multi-stage)
├── docker-compose.yml           # Декларативный запуск сервиса с лимитами cgroups
├── Makefile                     # Автоматизация типовых команд (run, test, race, lint)
├── go.mod                       # Описание Go-модуля и зависимостей
├── go.sum                       # Контрольные суммы используемых библиотек
└── README.md                    # Документация проекта
```

---

##  Быстрый старт (Старт / Стоп)

### Через Docker Compose
```bash
# Сборка и запуск контейнера:
make compose-up
# или напрямую: docker compose up --build

# Запуск в фоновом режиме:
docker compose up -d --build

# Остановка и удаление контейнеров:
make compose-down
# или напрямую: docker compose down
```

### Локальный запуск
```bash
# Запуск через Makefile:
make run

```

---

## Запуск тестов

```bash
# Обычные unit-тесты:
make test
# под капотом: go test -v ./...

# Тестирование с детектором состояния гонки данных (Data Race):
make test-race
# под капотом: go test -v -race -count=1 ./...

# Бенчмарки производительности ядра:
make bench
# под капотом: go test -bench=. -benchmem ./internal/broker/...

```

---

## Команда проекта

* **SKIND0A (Сергей)** — Архитектура ядра брокера, неблокирующая рассылка (Fan-Out), стресс-тесты под `-race`, бенчмарки, `Dockerfile`, `docker-compose.yml`, `Makefile`.
* **VladimirUlAeV (Владимир)** — Сетевой слой REST API и стриминг Server-Sent Events (`http.Flusher`), валидация DTO, сериализация JSON, unit-тесты транспорта.
