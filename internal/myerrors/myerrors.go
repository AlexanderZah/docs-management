package myerrors

import "errors"

var (
	ErrSessionNotFound = errors.New("session not found")
)
