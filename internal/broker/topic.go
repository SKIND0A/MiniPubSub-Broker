package broker

import (
	"sync"

	internal "github.com/SKIND0A/MiniPubSub/internal/error"

	"github.com/SKIND0A/MiniPubSub/internal/models"
)

type Topic struct {
	name string
	mu   sync.RWMutex
	subs map[string]*Subscriber
}

func NewTopic(name string) *Topic {
	return &Topic{
		name: name,
		subs: make(map[string]*Subscriber),
	}
}

// Закрываем топик
func (t *Topic) CloseTopic() {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, sub := range t.subs {
		sub.CloseChan()
	}
	clear(t.subs)
}

// Добавляем саба в топик
func (t *Topic) Subscribe(sub *Subscriber) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, ok := t.subs[sub.id]
	if ok {
		return internal.ErrSubscriberAlreadyExists
	}

	t.subs[sub.id] = sub
	return nil
}

// Удаляем саба из топика
func (t *Topic) Unsubscribe(id string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	sub, ok := t.subs[id]
	if !ok {
		return internal.ErrSubscriberNotFound
	}

	delete(t.subs, id)
	sub.CloseChan()
	return nil
}

// Рассылаем каждому сабу сообщение в его буфер канал
func (t *Topic) PublishMsg(msg models.Message) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, sub := range t.subs {
		select {
		case sub.ch <- msg:

		default:
			dropped++
		}
	}
	return dropped
}

