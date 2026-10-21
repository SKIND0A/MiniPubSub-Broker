package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"github.com/SKIND0A/MiniPubSub/internal/broker"
	internal "github.com/SKIND0A/MiniPubSub/internal/error"
	"github.com/SKIND0A/MiniPubSub/internal/models"
	"github.com/google/uuid"
)

type Handler struct{
	msg models.Message
	Broker *broker.Broker
}

func NewHandler(i models.Message, b *broker.Broker) *Handler{
	return &Handler{
		msg: i,
		Broker: b,
	}
}

func (h *Handler)NewPublishHandler(w http.ResponseWriter, r * http.Request){
	fmt.Println("Поступил запрос на публикацию")
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if r.Method != http.MethodPost{
		http.Error(w, internal.ErrRequstMethod.Error() ,http.StatusBadRequest)
		return
	}
	
	topic := r.PathValue("topic")
	if strings.TrimSpace(topic) == ""{
		http.Error(w, internal.ErrRequstTopic.Error(), http.StatusBadRequest)
		return
	}

	var err	error 
	var m models.Message
	var data models.MRequest
	userType := r.Header.Get("Content-type") //проверка контента
	if strings.HasPrefix(userType, "application/json"){
		unmarsh := json.NewDecoder(r.Body)
		if err := unmarsh.Decode(&data); err != nil{
			var syntaxError *json.SyntaxError
			var unmarshallTypeError *json.UnmarshalTypeError
			switch {
				case errors.As(err, &syntaxError):
					http.Error(w, internal.ErrRequstSyntax.Error(), http.StatusBadRequest)
					return
				case errors.As(err, &unmarshallTypeError):
					http.Error(w, internal.ErrUnmarshallRequst.Error(), http.StatusBadRequest)
					return
				case errors.Is(err, io.EOF): //пустое тело
					http.Error(w, internal.ErrBodyEmpty.Error(), http.StatusBadRequest)
					return
				case strings.HasPrefix(err.Error(), "json: unknown field"):
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				default:
					http.Error(w, internal.UnknownErr.Error(), http.StatusBadRequest)
					return
			}
		}
		m, err  = h.msg.NewMessage([]byte(data.PayLoad))
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}else if strings.HasPrefix(userType, "text/plain"){
		rawBody, err := io.ReadAll(r.Body)
		if err != nil{
			return
		}
		defer r.Body.Close()
		m, err = h.msg.NewMessage(rawBody)
		if err != nil{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}else{
		http.Error(w, internal.ErrRequstType.Error(), http.StatusUnsupportedMediaType)
		return
	}

	if err := h.Broker.PublishTopic(topic, m); err != nil{
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type","text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Успешно опубликовано, ID сообщения:" + m.ID))
}

func (h *Handler)AlltopicHandler(w http.ResponseWriter, r * http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, internal.ErrRequstMethod.Error(), http.StatusBadRequest)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 0)

	AllTopic, err := h.Broker.GetAllTopics()
	if err != nil{
		http.Error(w, err.Error(), http.StatusNoContent)
	}
	
	Output, err := json.Marshal(AllTopic)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusAccepted)
	w.Write(Output)

}

func (h *Handler)CreateTopicHandler(w http.ResponseWriter, r * http.Request){
	fmt.Println("Получен запрос на создание нового топика(темы)")
	if r.Method != http.MethodPost{
		http.Error(w, internal.ErrRequstMethod.Error(), http.StatusBadRequest)
		return
	}

	topic := r.PathValue("topic")
	if strings.TrimSpace(topic) == ""{
		http.Error(w, internal.ErrRequstTopic.Error(), http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 0)
	if _, err := h.Broker.CreateTopic(topic); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type","text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Создан топик " + topic))
}

func (h *Handler)SubscribeHandler(w http.ResponseWriter, r * http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, internal.ErrRequstMethod.Error(), http.StatusBadRequest)
		return
	}
	
	topic := r.PathValue("topic")
	if strings.TrimSpace(topic) == ""{
		http.Error(w, internal.ErrRequstTopic.Error(), http.StatusBadRequest)
		return
	}	

	flusher, ok := w.(http.Flusher)
	if ok == false{
		http.Error(w, internal.ErrRequstFlusher.Error(), http.StatusInternalServerError)
		return
	} 
	
	FindTopic, err := h.Broker.GetTopic(topic)
	if err != nil{
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	sub := broker.NewSubscriber(uuid.NewString(),100)
	FindTopic.Subscribe(sub)
	defer FindTopic.Unsubscribe(sub.ID())
	fmt.Println("Зарегестрирован новый подписчик")
	
	w.Header().Set("Content-type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	for {
		select{
			case msg := <-sub.Message() :
				responsemsg := models.ResponseMessage{
					ID: msg.ID,
					Body: string(msg.Body),
					TimeAt: msg.TimeAt,
				}
				OutputM, err := json.Marshal(responsemsg)
				if err != nil{
					continue
				}
				fmt.Fprintf(w,"data: %s\n\n", OutputM)
				flusher.Flush()
			case <-sub.Quit():
				return
			case <-r.Context().Done():
				return
		}
	}
}

func StartHTTPServer(h *Handler){
	defer h.Broker.CloseBroker()
	fmt.Println("Запуск HTTP-сервера, для остановки нажмите ctrl+C")
	fmt.Println("localhost:8080")
	mux := http.NewServeMux()
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	mux.HandleFunc("/api/v1/topics/{topic}/subscribe",h.SubscribeHandler)
	mux.HandleFunc("/api/v1/topics", h.AlltopicHandler)
	mux.HandleFunc("/api/v1/topics/{topic}",h.CreateTopicHandler)
	mux.HandleFunc("/api/v1/topics/{topic}/publish", h.NewPublishHandler)

		if err := server.ListenAndServe(); err != nil{
			if err.Error() == "http: Server closed"{
				fmt.Println(err.Error())
			}else{
				fmt.Print("Ошибка запуска по причине", err)
			}
		}
	
}