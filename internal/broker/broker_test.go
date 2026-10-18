package broker_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/SKIND0A/MiniPubSub/internal/broker"
	internal "github.com/SKIND0A/MiniPubSub/internal/error"
	"github.com/SKIND0A/MiniPubSub/internal/models"
)

func TestBroker_CreateTopic(t *testing.T) {
	b := broker.NewBroker()
	tests := []struct {
		name      string
		topicName string
		wantErr   bool
		targetErr error
	}{
		{
			name:      "успешное создание нового топика",
			topicName: "events",
			wantErr:   false,
			targetErr: nil,
		},
		{
			name:      "ошибка при создании дубликата",
			topicName: "events",
			wantErr:   true,
			targetErr: internal.ErrTopicAlreadyExists,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			topic, err := b.CreateTopic(test.topicName)

			if test.wantErr {
				if err == nil {
					t.Fatalf("ожидалась ошибка, получен nil")
				}
				if test.targetErr != nil && !errors.Is(err, test.targetErr) {
					t.Errorf("получена ошибка %v, ожидалась %v", err, test.targetErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if topic == nil {
				t.Fatalf("топик не должен быть nil")
			}
			savedTopic, err := b.GetTopic(test.topicName)
			if err != nil {
				t.Fatalf("не удалось получить созданный топик: %v", err)
			}
			if savedTopic != topic {
				t.Errorf("полученный топик не совпадает с созданным")
			}
		})
	}
}

func TestBroker_Subscribe_Unsubscribe(t *testing.T) {
	b := broker.NewBroker()
	topic, err := b.CreateTopic("work")
	if err != nil {
		t.Fatalf("Не удалось создать топик: %v", err)
	}

	sub := broker.NewSubscriber("sub#1", 100)

	if err := topic.Subscribe(sub); err != nil {
		t.Fatalf("ошибка успешной подписки: %v", err)
	}

	err = topic.Subscribe(sub)
	if !errors.Is(err, internal.ErrSubscriberAlreadyExists) {
		t.Errorf("ожидалась ошибка %v, получено %v", internal.ErrSubscriberAlreadyExists, err)
	}

	if err := topic.Unsubscribe(sub.ID()); err != nil {
		t.Fatalf("ошибка при отписке подписчика: %v", err)
	}
	err = topic.Unsubscribe(sub.ID())
	if !errors.Is(err, internal.ErrSubscriberNotFound) {
		t.Fatalf("ожидалась ошибка %v, получено %v", internal.ErrSubscriberNotFound, err)
	}

	select {
	case <-sub.Quit():
	default:
		t.Errorf("сигнальный канал подписчика не был закрыт после Unsubscribe")
	}
}

func TestBroker_Publish_FanOut_Drop(t *testing.T) {
	b := broker.NewBroker()
	topic, err := b.CreateTopic("work")
	if err != nil {
		t.Fatalf("Не удалось создать топик: %v", err)
	}

	subFast := broker.NewSubscriber("sub#1", 10)

	subSlow := broker.NewSubscriber("sub#2", 1)

	if err := topic.Subscribe(subFast); err != nil {
		t.Fatalf("ошибка подписки fastSub: %v", err)
	}
	if err := topic.Subscribe(subSlow); err != nil {
		t.Fatalf("ошибка подписки slowSub: %v", err)
	}

	msg1 := models.Message{
		ID:     "msg-1",
		Body:   []byte("первое тестовое сообщение"),
		TimeAt: time.Now(),
	}

	msg2 := models.Message{
		ID:     "msg-2",
		Body:   []byte("второе тестовое сообщение"),
		TimeAt: time.Now(),
	}

	msg3 := models.Message{
		ID:     "msg-3",
		Body:   []byte("третье тестовое сообщение"),
		TimeAt: time.Now(),
	}

	err = b.PublishTopic("work", msg1)
	if err != nil {
		t.Fatalf("не удалось опубликовать сообщение: %v", err)
	}
	err = b.PublishTopic("work", msg2)
	if err != nil {
		t.Fatalf("не удалось опубликовать сообщение: %v", err)
	}
	err = b.PublishTopic("work", msg3)
	if err != nil {
		t.Fatalf("не удалось опубликовать сообщение: %v", err)
	}

	for _, expectedID := range []string{"msg-1", "msg-2", "msg-3"} {
		select {
		case msg := <-subFast.Message():
			if msg.ID != expectedID {
				t.Errorf("fastSub: ожидали ID %s, получили %s", expectedID, msg.ID)
			}
		case <-time.After(50 * time.Millisecond):
			t.Fatalf("fastSub: таймаут ожидания сообщения %s", expectedID)
		}
	}

	select {
	case msg := <-subSlow.Message():
		if msg.ID != "msg-1" {
			t.Errorf("slowSub: первое сообщение должно быть msg-1, пришло %s", msg.ID)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatal("slowSub: таймаут ожидания первого сообщения")
	}

	select {
	case msg := <-subSlow.Message():
		t.Errorf("slowSub: буфер должен быть пуст, но получено неожиданное сообщение: %s", msg.ID)
	default:

	}
	if count := b.PublishedCount(); count != 3 {
		t.Errorf("ожидалось 3 опубликованных сообщения в метриках, получено: %d", count)
	}
}

func TestBroker_ConcurrentPublish_Race(t *testing.T) {
	b := broker.NewBroker()
	topicName := "stress-topic"

	topic, err := b.CreateTopic(topicName)
	if err != nil {
		t.Fatalf("не удалось создать топик: %v", err)
	}

	const (
		numSubs           = 5
		numPublishers     = 50
		msgsPerPublisher  = 50
		expectedPublished = numPublishers * msgsPerPublisher
	)

	subs := make([]*broker.Subscriber, numSubs)
	for i := 0; i < numSubs; i++ {
		subID := fmt.Sprintf("sub-%d", i)
		sub := broker.NewSubscriber(subID, 100)
		if err := topic.Subscribe(sub); err != nil {
			t.Fatalf("ошибка подписки %s: %v", subID, err)
		}
		subs[i] = sub
	}

	var subWg sync.WaitGroup
	subWg.Add(numSubs)

	for _, sub := range subs {
		go func(s *broker.Subscriber) {
			defer subWg.Done()
			for {
				select {
				case _, ok := <-s.Message():
					if !ok {
						return
					}

				case <-s.Quit():

					return
				}
			}
		}(sub)
	}

	var pubWg sync.WaitGroup
	pubWg.Add(numPublishers)

	for p := 0; p < numPublishers; p++ {
		go func(pubID int) {
			defer pubWg.Done()
			for m := 0; m < msgsPerPublisher; m++ {
				msg := models.Message{
					ID:     fmt.Sprintf("pub-%d-msg-%d", pubID, m),
					Body:   []byte("stress-test-payload"),
					TimeAt: time.Now(),
				}
				if err := b.PublishTopic(topicName, msg); err != nil {
					t.Errorf("ошибка публикации сообщения: %v", err)
				}
			}
		}(p)
	}

	pubWg.Wait()

	time.Sleep(50 * time.Millisecond)

	b.CloseBroker()

	subWg.Wait()

	if count := b.PublishedCount(); count != uint64(expectedPublished) {
		t.Errorf("несовпадение счетчика метрик: ожидалось %d, получено %d", expectedPublished, count)
	}
}

func TestBroker_Close_RejectsNewOperations(t *testing.T) {
	b := broker.NewBroker()
	topicName := "shutdown-topic"

	topic, err := b.CreateTopic(topicName)
	if err != nil {
		t.Fatalf("не удалось создать топик: %v", err)
	}

	sub := broker.NewSubscriber("sub-test", 10)
	if err := topic.Subscribe(sub); err != nil {
		t.Fatalf("не удалось подписать клиента: %v", err)
	}

	
	b.CloseBroker()

	b.CloseBroker()

	select {
	case <-sub.Quit():
	default:
		t.Errorf("сигнальный канал подписчика остался открытым после CloseBroker")
	}

	_, err = b.CreateTopic("another-topic")
	if err == nil {
		t.Errorf("ожидалась ошибка при создании топика после CloseBroker, получен nil")
	} else if !errors.Is(err, internal.ErrTopicClosed) {
		t.Logf("получена ошибка создания топика: %v", err)
	}

	msg := models.Message{
		ID:     "msg-after-close",
		Body:   []byte("test"),
		TimeAt: time.Now(),
	}
	err = b.PublishTopic(topicName, msg)
	if err == nil {
		t.Errorf("ожидалась ошибка публикации после CloseBroker, получен nil")
	}

	newSub := broker.NewSubscriber("new-sub", 10)
	err = topic.Subscribe(newSub)
	if err == nil {
		t.Errorf("ожидалась ошибка при попытке подписки на закрытый топик, получен nil")
	}
}

func BenchmarkBroker_Publish(b *testing.B) {
	br := broker.NewBroker()
	topicName := "bench-topic"

	topic, err := br.CreateTopic(topicName)
	if err != nil {
		b.Fatalf("не удалось создать тему для бенчмарка: %v", err)
	}

	sub := broker.NewSubscriber("bench-sub", 2048)
	if err := topic.Subscribe(sub); err != nil {
		b.Fatalf("ошибка подписки: %v", err)
	}

	go func() {
		for {
			select {
			case _, ok := <-sub.Message():
				if !ok {
					return
				}
			case <-sub.Quit():
				return
			}
		}
	}()

	msg := models.Message{
		ID:     "bench-msg",
		Body:   []byte("benchmark-payload-bytes"),
		TimeAt: time.Now(),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = br.PublishTopic(topicName, msg)
	}

	b.StopTimer()
	br.CloseBroker()
}