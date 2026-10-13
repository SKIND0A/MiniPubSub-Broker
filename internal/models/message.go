package models

import (
	"time"
)

type Message struct {
	ID     string `json:"id"`
	Body   []byte `json:"body"`
	TimeAt time.Time `json:"timeAt"`
}

type MRequest struct {
	PayLoad []byte `json:"payload"`
}
