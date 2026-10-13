package broker

import "github.com/SKIND0A/MiniPubSub/internal/models"

type Subscriber struct {
	id   string
	ch   chan models.Message 
	quit chan struct{}
}

func (s *Subscriber) Medssage() <-chan models.Message {
	return s.ch
}

// Закрытие буф канала Unsubscribe()
func (s *Subscriber) CloseChan() {
	close(s.quit)
}
