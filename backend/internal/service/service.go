package service

import "sync"

type (
	AuthToken string
	Email     string
	Koinu     int64
)

type User struct {
	Email Email
}

type Service struct {
	// don't do this in a real service, obviously 😄
	sessions      map[AuthToken]*User
	sessionsMutex sync.RWMutex

	// don't do this in a real service, obviously 😄
	balances      map[Email]Koinu
	balancesMutex sync.RWMutex
}

func New() *Service {
	return &Service{
		sessions: make(map[AuthToken]*User),
		balances: map[Email]Koinu{
			"void@example.com":  1_234_567_890,
			"shiba@example.com": 2_111_222_333,
		},
	}
}

func (s *Service) Close() {
}
