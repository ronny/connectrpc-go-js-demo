package service

import (
	"context"
	"crypto/rand"
)

type LoginRequest struct {
	Email    Email
	Password string
}

type LoginResponse struct {
	AuthToken AuthToken
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	s.sessionsMutex.Lock()
	defer s.sessionsMutex.Unlock()

	if req.Password != "dogecoin" {
		return nil, ErrUnauthenticated
	}

	user := &User{Email: req.Email}
	token := AuthToken(rand.Text())
	s.sessions[token] = user

	return &LoginResponse{AuthToken: token}, nil
}
