package model

import "errors"

var (
	ErrEmptyString = errors.New("empty string can't be encrypted")
)

// TODO: Newtypes over slices are auto casting if placed in separate module
// i think its due to fat pointer
type Encrypted[T any] []byte

func (e Encrypted[T]) String() string {
	return string(e)
}

func EncryptedFromString[T any](s string) (Encrypted[T], error) {
	if s == "" {
		return Encrypted[T]{}, ErrEmptyString
	}
	return Encrypted[T](s), nil
}
