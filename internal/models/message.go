package models

import (
	"bytes"
	"time"

	internal "github.com/SKIND0A/MiniPubSub/internal/error"
	"github.com/google/uuid"
)

type Message struct {
	ID     string 
	Body   []byte 
	TimeAt time.Time 
}

type ResponseMessage struct {
	ID     string `json:"id"`
	Body   string `json:"body"`
	TimeAt time.Time `json:"timeAt"`
}

type MRequest struct {
	PayLoad string `json:"payload"`
}

func (s *Message) NewMessage(body []byte) (Message, error){
	if len(bytes.TrimSpace(body)) == 0{
		return Message{}, internal.ErrBodyEmpty
	}

	s.ID = uuid.NewString()
	s.Body = body
	s.TimeAt = time.Now()

	return *s, nil
}
