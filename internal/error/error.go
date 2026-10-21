package internal

import "errors"

var (
	//sub errors
	ErrSubscriberAlreadyExists = errors.New("Подписчик с таким id уже существует")
	ErrSubscriberNotFound      = errors.New("Подписчик с таким id не найден")

	//topic errors
	ErrTopicAlreadyExists = errors.New("Топик с таким именем уже существует")
	ErrTopicNotFound      = errors.New("Топик с таким именем не найден")
	ErrTopicClosed        = errors.New("Топик закрыт и не может быть открыт")
	ErrTopicNotCreated    = errors.New("Топики еще не созданы")

	//broker errors
	ErrBrokerClosed = errors.New("Брокер закрыт и не может быть создан")

	// DTO error
	ErrBodyEmpty = errors.New("Пустое тело")

	//http error
	ErrRequstTopic      = errors.New("Не указана тема(топик)")
	ErrRequstMethod     = errors.New("Используется не подходящий метод запроса")
	ErrRequstType       = errors.New("Используется не подходящий тип данных в запросе")
	ErrRequstSyntax     = errors.New("Ошибка синтаксиса в запросе")
	ErrUnmarshallRequst = errors.New("Ошибка анмаршлинга запроса")
	UnknownErr          = errors.New("Неизвестная ошибка")
	ErrRequstFlusher    = errors.New("Запрос не может работать в режиме стриминга")
)
