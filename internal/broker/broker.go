package broker

import (
	"sync"
	"sync/atomic"

	internal "github.com/SKIND0A/MiniPubSub/internal/error"
	"github.com/SKIND0A/MiniPubSub/internal/models"
)

type Broker struct {
	mu      sync.RWMutex
	topics  map[string]*Topic
	counter atomic.Uint64
}

func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]*Topic),
	}
}

func (b *Broker) PublishedCount() uint64 {
	return b.counter.Load()
}

// Создаем новый топик
func (b *Broker) CreateTopic(name string) (*Topic, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, ok := b.topics[name]
	if ok {
		return nil, internal.ErrTopicAlreadyExists
	}

	topic := NewTopic(name)
	b.topics[name] = topic

	return topic, nil
}

// Получаем топик
func (b *Broker) GetTopic(name string) (*Topic, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topic, ok := b.topics[name]
	if !ok {
		return nil, internal.ErrTopicNotFound
	}
	return topic, nil
}

// Отправка топика
func (b *Broker) PublishTopic(topicName string, msg models.Message) error {

	topic, err := b.GetTopic(topicName)
	if err != nil {
		return err
	}
	topic.PublishMsg(msg)
	b.counter.Add(1)
	return nil
}

// Закрытие топика
func (b *Broker) CloseBroker() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, topic := range b.topics {
		topic.CloseTopic()
	}
	clear(b.topics)
}
