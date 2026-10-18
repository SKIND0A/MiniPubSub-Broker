package models

import (
	"time"
	"github.com/google/uuid"
	internal "github.com/SKIND0A/MiniPubSub/internal/error"
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

func (s *Message) NewMessage(body []byte) (error, Message){
	if body == nil{
		return internal.ErrBodyEmpty, Message{}
	}

	s.ID = uuid.NewString()
	s.Body = body
	s.TimeAt = time.Now()
	
	return nil, *s
}
