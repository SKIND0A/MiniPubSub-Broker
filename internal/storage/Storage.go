package storage

import (
	"sync"
	"time"

	"github.com/SKIND0A/MiniPubSub/internal/models"
	"github.com/google/uuid"
)

type StorageMessage struct {
	messagers map[string]models.Message
	mu        sync.Mutex
}

func NewStorageM() *StorageMessage {
	return &StorageMessage{}
}

func (s *StorageMessage) NewMessage(body []byte) {
	Id := uuid.NewString()
	t := time.Now()

	s.messagers[Id] = models.Message{
		ID: Id,
		Body: body,
		TimeAt: t,

	}
	
}

func (s *StorageMessage) DeleteMessage(Id string)error{
	if Id == ""{
		return //error
	}
	if _ ,ok := s.messagers[Id]; ok == true{
		return //error
	}
	delete(s.messagers, Id)
	return nil
}
