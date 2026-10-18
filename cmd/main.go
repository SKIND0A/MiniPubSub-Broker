package main

import (

	"github.com/SKIND0A/MiniPubSub/internal/broker"
	"github.com/SKIND0A/MiniPubSub/internal/handler"
	"github.com/SKIND0A/MiniPubSub/internal/models"
)

func main() {
	Broker := broker.NewBroker()
	Handler := handler.NewHandler(models.Message{}, Broker)
	handler.StartHTTPServer(Handler)
}
