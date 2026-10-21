package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SKIND0A/MiniPubSub/internal/broker"
	"github.com/SKIND0A/MiniPubSub/internal/handler"
	"github.com/SKIND0A/MiniPubSub/internal/models"
)

func TestNewPublishHandler(t *testing.T) {
	t.Run("запрос с неправельным методом", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()

		req_method := httptest.NewRequest("GET", "/api/v1/topics/Тест/publish", strings.NewReader(""))
		h.NewPublishHandler(rec, req_method)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("не прошла валидация на метод, текущий код: %d", rec.Code)
		}
	})

	t.Run(("запрос на не существующий топик"), func(T *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()

		req_method := httptest.NewRequest("POST", "/api/v1/topics/Тест/publish", strings.NewReader("111"))
		h.NewPublishHandler(rec, req_method)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("не прошла валидация на пустой топик, текущий код: %d", rec.Code)
		}
	})

	t.Run("запрос текстом", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()
		if _, err := h.Broker.CreateTopic("Тест"); err != nil {
			t.Fatalf("%s", "Ошибка создания топика по причине:"+err.Error())
			return
		}

		req_method := httptest.NewRequest("POST", "/api/v1/topics/Тест/publish", strings.NewReader("111"))
		req_method.Header.Set("Content-Type", "text/plain")
		req_method.SetPathValue("topic", "Тест")

		h.NewPublishHandler(rec, req_method)
		if rec.Code != http.StatusCreated {
			t.Errorf("не прошла валидация на формат данных в виде текста текущий код: %d", rec.Code)
		}

		frec := httptest.NewRecorder()
		freq_method := httptest.NewRequest("POST", "/api/v1/topics/Тест/publish", strings.NewReader(""))
		freq_method.Header.Set("Content-Type", "text/plain")
		freq_method.SetPathValue("topic", "Тест")

		h.NewPublishHandler(frec, freq_method)
		if frec.Code != http.StatusBadRequest {
			t.Errorf("прошла валидация на формат данных в виде пустого текста текущий код: %d", frec.Code)
		}
	})

	t.Run("запрос json", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()
		if _, err := h.Broker.CreateTopic("Тест"); err != nil {
			t.Fatalf("%s", "Ошибка создания топика по причине:"+err.Error())
		}

		type structTest struct {
			Body string `json:"payload"`
		}
		data, err := json.Marshal(structTest{
			Body: "Тест",
		})
		if err != nil {
			t.Fatalf("%s", "Ошибка маршалинга по причине:"+err.Error())
			return
		}

		fdata, err := json.Marshal(structTest{
			Body: "",
		})
		if err != nil {
			t.Fatalf("%s", "Ошибка маршалинга по причине:"+err.Error())
			return
		}

		req_method := httptest.NewRequest("POST", "/api/v1/topics/Тест/publish", bytes.NewReader(data))
		req_method.Header.Set("Content-Type", "application/json")
		req_method.SetPathValue("topic", "Тест")

		h.NewPublishHandler(rec, req_method)
		if rec.Code != http.StatusCreated {
			t.Errorf("не прошла валидация на формат данных в виде json, текущий код: %d", rec.Code)
		}

		frec := httptest.NewRecorder()
		freq_method := httptest.NewRequest("POST", "/api/v1/topics/Тест/publish", bytes.NewReader(fdata))
		freq_method.Header.Set("Content-Type", "application/json")
		freq_method.SetPathValue("topic", "Тест")

		h.NewPublishHandler(frec, freq_method)
		if frec.Code != http.StatusBadRequest {
			t.Errorf("прошла валидация на формат данных в виде пустого json, текущий код: %d", frec.Code)
		}
	})

	t.Run("запрос с неправельным форматом данных", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()

		req_method := httptest.NewRequest("POST", "/api/v1/topics/Тест/publish", nil)
		req_method.SetPathValue("topic", "Тест")

		h.NewPublishHandler(rec, req_method)
		if rec.Code != http.StatusUnsupportedMediaType {
			t.Errorf("прошла валидация на формат данных в виде , текущий код: %d", rec.Code)
		}

	})
}

func TestAlltopicHandler(t *testing.T) {
	t.Run("запрос с неправельным методом", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()

		req_method := httptest.NewRequest("POST", "/api/v1/topics", strings.NewReader(""))
		h.AlltopicHandler(rec, req_method)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("не прошла валидация на метод, текущий код: %d", rec.Code)
		}
	})

	t.Run("запрос на пустые топики", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()

		req_method := httptest.NewRequest("GET", "/api/v1/topics", strings.NewReader("111"))
		h.AlltopicHandler(rec, req_method)
		if rec.Code != http.StatusNoContent {
			t.Errorf("не прошла валидация на метод, текущий код: %d", rec.Code)
		}
	})

	t.Run("запрос на выдачу топиков", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		if _, err := h.Broker.CreateTopic("Тест"); err != nil {
			t.Fatalf("%s", "Ошибка создания топика по причине:"+err.Error())
			return
		}

		rec := httptest.NewRecorder()

		req_method := httptest.NewRequest("GET", "/api/v1/topics", strings.NewReader(""))
		h.AlltopicHandler(rec, req_method)
		if rec.Code != http.StatusAccepted {
			t.Errorf("не прошла валидация на тело, текущий вывод: %s", rec.Body.String())
		}
	})
}

func TestCreateTopicHandler(t *testing.T) {
	t.Run("запрос с неправельным методом", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()

		req_method := httptest.NewRequest("GET", "/api/v1/topics/Тест", strings.NewReader(""))
		h.CreateTopicHandler(rec, req_method)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("не прошла валидация на метод, текущий код: %d", rec.Code)
		}
	})

	t.Run("запрос с пустым топиком", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()

		req_method := httptest.NewRequest("POST", "/api/v1/topics/", strings.NewReader(""))
		h.CreateTopicHandler(rec, req_method)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("не прошла валидация на метод, текущий код: %d", rec.Code)
		}
	})

	t.Run("запрос на создание топика", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()

		req_method := httptest.NewRequest("POST", "/api/v1/topics/Тест", nil)
		req_method.SetPathValue("topic","Тест")
		h.CreateTopicHandler(rec, req_method)
		if rec.Code != http.StatusCreated {
			t.Errorf("не прошла валидация на создание, текущий код: %d", rec.Code)
		}
	})
}

func TestSubscribeHandler(t *testing.T) {
	t.Run("запрос с неправельным методом", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()
		if _, err := h.Broker.CreateTopic("Тест"); err != nil {
			t.Fatalf("%s", err.Error())
		}

		req := httptest.NewRequest("POST", "/api/v1/topics/Тест/subscribe", strings.NewReader(""))
		h.SubscribeHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("валидация не пройдена, текущий код: %d", rec.Code)
		}
	})

	t.Run("запрос с неправельным методом", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()
		if _, err := h.Broker.CreateTopic("Тест"); err != nil {
			t.Fatalf("%s", err.Error())
		}

		req := httptest.NewRequest("POST", "/api/v1/topics/Тест/subscribe", strings.NewReader(""))
		h.SubscribeHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("валидация не пройдена, текущий код: %d", rec.Code)
		}
	})

	t.Run("запрос на пустой топик", func(t *testing.T) {
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()
		if _, err := h.Broker.CreateTopic("Тест"); err != nil {
			t.Fatalf("%s", err.Error())
		}

		req := httptest.NewRequest("GET", "/api/v1/topics/Тест/subscribe", strings.NewReader(""))
		req.SetPathValue("topic", "")
		h.SubscribeHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("валидация не пройдена, текущий код: %d", rec.Code)
		}
	})

	t.Run("запрос на подписку", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan struct{})
		b := broker.NewBroker()
		h := handler.NewHandler(models.Message{}, b)
		rec := httptest.NewRecorder()
		if _, err := h.Broker.CreateTopic("Тест"); err != nil {
			t.Fatalf("%s", err.Error())
		}

		req := httptest.NewRequest("GET", "/api/v1/topics/Тест/subscribe", nil).WithContext(ctx)
		req.SetPathValue("topic", "Тест")
		go func() {
			h.SubscribeHandler(rec, req)
			close(done)
		}()

		time.Sleep(1 * time.Second)
		rech := httptest.NewRecorder()
		reqh := httptest.NewRequest("POST", "/api/v1/topics/Тест/publish", strings.NewReader("Тест"))
		reqh.Header.Set("Content-Type", "text/plain")
		reqh.SetPathValue("topic", "Тест")
		h.NewPublishHandler(rech, reqh)

		cancel()
		select {
		case <-done:
			if rec.Body.String() == "" {
				t.Errorf("валидация не пройдена %s", rec.Body.String())
				return
			}
		case <-time.After(1 * time.Second):
			t.Fatal("Хендлер завис после отмены")
		}

	})
}
