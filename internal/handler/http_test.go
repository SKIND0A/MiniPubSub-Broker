package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SKIND0A/MiniPubSub/internal/broker"
	"github.com/SKIND0A/MiniPubSub/internal/handler"
	"github.com/SKIND0A/MiniPubSub/internal/models"
)

func TestNewPublishHandler(t *testing.T){
	req := httptest.NewRequest("POST","/api/v1/topics/Тест/publish", strings.NewReader("111"))
	rec := httptest.NewRecorder()

	b := broker.NewBroker()
	h := handler.NewHandler(models.Message{}, b)
	h.Broker.CreateTopic("Тест")
	h.NewPublishHandler(rec, req)
	if rec.Code != http.StatusBadRequest{
		t.Errorf("ждали статус 400, а сервер вернул %d", rec.Code)
	}

}