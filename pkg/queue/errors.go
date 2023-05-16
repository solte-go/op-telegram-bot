package queue

import "errors"

var (
	contextExceeded                 = errors.New("context exceeded")
	ErrNoConnectionWithProvidedName = errors.New("no connection with provided name")
	ErrEmptyConnectionName          = errors.New("empty connection name not supported")
	ErrUnknownMessageType           = errors.New("unknown message type")
)
