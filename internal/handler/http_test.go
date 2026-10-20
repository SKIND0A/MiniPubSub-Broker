package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SKIND0A/MiniPubSub/internal/broker"
	"github.com/SKIND0A/MiniPubSub/internal/handler"
	"github.com/SKIND0A/MiniPubSub/internal/models"
)

func TestNewPublishHandler(t *testing.T){
	t.Run("запрос с неправельным методом", func(t *testing.T){
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()

		req_method_fail := httptest.NewRequest("GET","/api/v1/topics/Тест/publish", strings.NewReader(""))
		h.NewPublishHandler(rec, req_method_fail)
		if rec.Code != http.StatusBadRequest{
			t.Errorf("не прошла валидация на метод, код: %d", rec.Code)
		}
	})
	
	t.Run("запрос текстом", func(t *testing.T){
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()
		if _ ,err := h.Broker.CreateTopic("Тест"); err != nil{
			return
		}

		req_method_fail := httptest.NewRequest("POST","/api/v1/topics/Тест/publish", strings.NewReader("111"))
		req_method_fail.Header.Set("Content-Type", "text/plain")
		req_method_fail.SetPathValue("topic", "Тест")
		
		h.NewPublishHandler(rec, req_method_fail)
		if rec.Code != http.StatusCreated{
			t.Errorf("не прошла валидация на формат данных в виде текста%d", rec.Code)
		}
	})

	t.Run("запрос json", func(t * testing.T){
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()
		if _ ,err := h.Broker.CreateTopic("Тест"); err != nil{
			fmt.Println(err)
		}

		type structTest struct{
			Body string `json:"body"`
		}
		data, err := json.Marshal(structTest{
			Body: "Тест",
		})
		if err != nil{
			return
		}

		req_method_fail := httptest.NewRequest("POST","/api/v1/topics/Тест/publish", bytes.NewReader(data))
		req_method_fail.Header.Set("Content-Type", "application/json")
		req_method_fail.SetPathValue("topic", "Тест")
		
		h.NewPublishHandler(rec, req_method_fail)
		if rec.Code != http.StatusCreated{
			t.Errorf("не прошла валидация на формат данных в виде json, код: %d", rec.Code)
		}
	})

	t.Run("запрос с неправельным форматом данных", func(t *testing.T){
		

	})
}