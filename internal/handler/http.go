package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Handler struct{}


func StartHTTPServer(){
	fmt.Println("Запуск HTTP-сервера, для остановки нажмите ctrl+C")
	fmt.Println("localhost:8080")
	mux := http.NewServeMux()
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	// handlers

	go func(){
		if err := server.ListenAndServe(); err != nil{
			//error prodaction
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