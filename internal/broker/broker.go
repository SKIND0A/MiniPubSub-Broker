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
	closed  bool
}

func NewBroker() *Broker {
	return  &Broker{
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

	if b.closed {
		return nil, internal.ErrTopicClosed
	}

	_, ok := b.topics[name]
	if ok {
		return nil, internal.ErrTopicAlreadyExists
	}

	topic := NewTopic(name)
	b.topics[name] = topic

	return topic, nil
}

// получаем все топики
func (b *Broker) GetAllTopics() ( []string, error){
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return  nil, internal.ErrTopicClosed
	}
	
	if len(b.topics) == 0{
		return nil, internal.ErrTopicNotCreated
	}

	AllTopics := make([]string, 0, len(b.topics))
	for name := range b.topics{
		AllTopics = append(AllTopics, name)
	}
	return AllTopics, nil
}
// Получаем топик
func (b *Broker) GetTopic(name string) (*Topic, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return nil, internal.ErrTopicClosed
	}

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
	b.counter.Add(1)
	topic.PublishMsg(msg)
	return nil
}

// Закрытие топика
func (b *Broker) CloseBroker() {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if b.closed {
		return
	}
	b.closed = true

	for _, topic := range b.topics {
		topic.CloseTopic()
	}
	clear(b.topics)
}
