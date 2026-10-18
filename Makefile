.PHONY: build run test test-race bench lint docker-build docker-run help

# Сборка бинарника
build:
	go build -o bin/pubsub cmd/main.go

# Сборка и локальный запуск приложения
run: build
	./bin/pubsub

# Запуск unit-тестов с подробным выводом
test:
	go test -v ./...

# Запуск тестов со встроенным детектором гонок данных и сбросом кэша
test-race:
	go test -v -race -count=1 ./...

# Запуск бенчмарка производительности с замером аллокаций памяти
bench:
	go test -bench=BenchmarkBroker_Publish -benchmem -run=^$$ ./internal/broker/...

# Сборка  Docker-образа
docker-build:
	docker build -t minipubsub:latest .

# Запуск собранного контейнера с пробросом порта
docker-run:
	docker run -p 8080:8080 --rm minipubsub:latest