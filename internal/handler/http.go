package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SKIND0A/MiniPubSub/internal/models"
)

type Handler struct{
	storage models.StorageMessage
}

func NewHandler(i models.StorageMessage) *Handler{
	return &Handler{
		storage: i,
	}
}

func (h *Handler)NewMessege(w http.ResponseWriter, r * http.Request){
	
}

func StartHTTPServer(h * Handler){
	fmt.Println("Запуск HTTP-сервера, для остановки нажмите ctrl+C")
	fmt.Println("localhost:8080")
	mux := http.NewServeMux()
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	mux.HandleFunc("/newmessege", h.NewMessege)

	go func(){
		if err := server.ListenAndServe(); err != nil{
			if err.Error() == "http: Server closed"{}else{
				fmt.Print("Ошибка запуска по причине", err)
			}
		}
	}()

	//Улавливать килл 
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil{
		// error prodaction
	}
	
}