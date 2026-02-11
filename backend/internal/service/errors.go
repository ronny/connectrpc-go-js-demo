package service

import "errors"

var (
	ErrUnauthenticated     = errors.New("unauthenticated")
	ErrInsufficientBalance = errors.New("insufficient balance")
)
