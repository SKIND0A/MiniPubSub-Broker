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
	closed bool
}

var once sync.Once

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
	if t.closed {
		return 
	}
	t.closed = true

	for _, sub := range t.subs {
		sub.CloseChan()
	}
	clear(t.subs)
}

// Добавляем саба в топик
func (t *Topic) Subscribe(sub *Subscriber) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return internal.ErrTopicClosed
	}

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

	if t.closed {
		return internal.ErrTopicClosed
	}

	sub, ok := t.subs[id]
	if !ok {
		return internal.ErrSubscriberNotFound
	}

	delete(t.subs, id)
	once.Do(sub.CloseChan) // однократное выполнение функции CloseChan
	return nil
}

// Рассылаем каждому сабу сообщение в его буфер канал
func (t *Topic) PublishMsg(msg models.Message) (int, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.closed {
		return 0, internal.ErrTopicClosed
	}

	dropped := 0
	for _, sub := range t.subs {
		select {
		case sub.ch <- msg:

		default:
			dropped++
		}
	}
	return dropped, nil
}
