package storage

import (
	"sync"
	"time"

	internal "github.com/SKIND0A/MiniPubSub/internal/error"
	"github.com/SKIND0A/MiniPubSub/internal/models"
	"github.com/google/uuid"
)

type StorageMessage struct {
	messagers map[string]models.Message
	mu        sync.Mutex
}

func NewStorageM() *StorageMessage {
	return &StorageMessage{messagers: make(map[string]models.Message)}
}

func (s *StorageMessage) NewMessage(body []byte) error{
	if body == nil{
		return internal.ErrBodyEmpty
	}

	Id := uuid.NewString()
	t := time.Now()
	s.mu.Lock()
	s.messagers[Id] = models.Message{
		ID: Id,
		Body: body,
		TimeAt: t,

	}
	s.mu.Unlock()
	return nil
}

func (s *StorageMessage) DeleteMessage(Id string)error{
	if Id == ""{
		return internal.ErrIDEmpty
	}

	s.mu.Lock()
	if _ ,ok := s.messagers[Id]; !ok{
		s.mu.Unlock()
		return internal.ErrIDNotFound
	}
	s.mu.Unlock()

	delete(s.messagers, Id)
	return nil
}
