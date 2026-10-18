FROM golang:1.26.4-alpine

WORKDIR /Broker-app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app cmd/main.go

RUN adduser -D -g '' appbroker
USER appbroker

EXPOSE 8080

ENTRYPOINT ["./app"]