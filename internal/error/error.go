package internal

import "errors"

var (
	ErrSubscriberAlreadyExists = errors.New("Подписчик с таким id уже существует")
	ErrSubscriberNotFound      = errors.New("Подписчик с таким id не найден")
	ErrTopicAlreadyExists      = errors.New("Топик с таким именем уже существует")
	ErrTopicNotFound           = errors.New("Топик с таким именем не найден")
)
