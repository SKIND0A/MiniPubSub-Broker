package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	internal "github.com/SKIND0A/MiniPubSub/internal/error"
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

func (h *Handler)NewPublish(w http.ResponseWriter, r * http.Request){
	fmt.Println("Поступил запрос на публикацию")

	if r.Method != http.MethodPost{
		http.Error(w, internal.ErrRequstMethod.Error() ,http.StatusBadRequest)
		return
	}
	
	topic := r.PathValue("topic")
	if strings.TrimSpace(topic) == ""{
		http.Error(w, internal.ErrRequstTopic.Error(), http.StatusBadRequest)
		return
	}

	userType := r.Header.Get("Content-type") //проверка контента
	if !strings.HasPrefix(userType, "application/json"){
		http.Error(w, internal.ErrRequstType.Error(), http.StatusUnsupportedMediaType)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	var data models.MRequest
	unmarsh := json.NewDecoder(r.Body)
	if err := unmarsh.Decode(&data); err != nil{
		var syntaxError *json.SyntaxError
		var unmarshallTypeError *json.UnmarshalTypeError
		switch {
			case errors.As(err, &syntaxError):
				http.Error(w, internal.ErrRequstSyntax.Error(), http.StatusBadRequest)
			case errors.As(err, &unmarshallTypeError):
				http.Error(w, internal.ErrUnmarshallRequst.Error(), http.StatusBadRequest)
			case errors.Is(err, io.EOF): //пустое тело
				http.Error(w, internal.ErrBodyEmpty.Error(), http.StatusBadRequest)
			case strings.HasPrefix(err.Error(), "json: unknown field"):
				http.Error(w, err.Error(), http.StatusBadRequest)
			default:
				http.Error(w, internal.UnknownErr.Error(), http.StatusBadRequest)
			}
	}

	if err := h.storage.NewMessage([]byte(data.PayLoad)); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	
	w.Header().Set("Content-Type","applictaion/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))

}

func StartHTTPServer(h * Handler){
	fmt.Println("Запуск HTTP-сервера, для остановки нажмите ctrl+C")
	fmt.Println("localhost:8080")
	mux := http.NewServeMux()
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	mux.HandleFunc("/api/v1/topics/{topic}/publish", h.NewPublish)

	go func(){
		if err := server.ListenAndServe(); err != nil{
			if err.Error() == "http: Server closed"{
				fmt.Println(err.Error())
			}else{
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