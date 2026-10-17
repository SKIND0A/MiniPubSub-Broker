package internal

import "errors"

var (
	ErrSubscriberAlreadyExists = errors.New("Подписчик с таким id уже существует")
	ErrSubscriberNotFound      = errors.New("Подписчик с таким id не найден")
	ErrTopicAlreadyExists      = errors.New("Топик с таким именем уже существует")
	ErrTopicNotFound           = errors.New("Топик с таким именем не найден")

	// DTO error
	ErrBodyEmpty = errors.New("Пустое тело")
	ErrIDEmpty = errors.New("Пустое айди")
	ErrIDNotFound = errors.New("Сообщение с таким айди не найдено")

	//http error
	ErrRequstTopic = errors.New("Не указана тема(топик)")
	ErrRequstMethod = errors.New("Используется не подходящий метод запроса")
	ErrRequstType = errors.New("Используется не подходящий тип данных в запросе")
	ErrRequstSyntax = errors.New("Ошибка синтаксиса в запросе")
	ErrUnmarshallRequst = errors.New("Ошибка анмаршлинга запроса")
	UnknownErr = errors.New("Неизвестная ошибка")
	ErrRequstFlusher = errors.New("Запрос не может работать в режиме стриминга")
)
