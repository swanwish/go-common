package cache

import "errors"

var (
	ErrKeyNotFound = errors.New("Key not found")
)

type Cache interface {
	SetValue(key, value []byte) error
	GetValue(key []byte) ([]byte, error)
	Delete(key []byte) error
}
