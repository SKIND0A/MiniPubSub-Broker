package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/SKIND0A/MiniPubSub/internal/broker"
	"github.com/SKIND0A/MiniPubSub/internal/handler"
	"github.com/SKIND0A/MiniPubSub/internal/models"
)

func main() {
	fmt.Println("CLI-управление Broker")
	Broker := broker.NewBroker()
	Handler := handler.NewHandler(models.Message{}, Broker)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		Input := scanner.Text()
		switch Input {
		case "1":
			handler.StartHTTPServer(Handler)
		case "2":
		case "3":
		case "0":
			fmt.Println("Выход из программы")
			return
		default:
			fmt.Println("Нету такой команды")
		}
	}

}
