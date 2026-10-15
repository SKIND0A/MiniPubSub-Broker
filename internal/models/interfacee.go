package models

type StorageMessage interface{
	NewMessage(body []byte) error
	DeleteMessage(Id string)error
}