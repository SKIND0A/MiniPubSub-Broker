package storage

import (
	"sync"
	"github.com/SKIND0A/MiniPubSub/internal/models"
)

type StorageMessage struct {
	messagers map[string]models.Message
	mu        sync.Mutex
}

func NewStorageM()*StorageMessage{
	return &StorageMessage{}
}



