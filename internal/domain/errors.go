package domain

import "errors"

var (
	ErrUrlNotFound  = errors.New("url not found")
	ErrTimeout      = errors.New("timeout")
	ErrInvalidInput = errors.New("invalid input")
)
