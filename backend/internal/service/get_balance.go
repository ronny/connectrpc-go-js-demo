package service

import (
	"context"
)

type GetBalanceRequest struct {
	Email Email
}

type GetBalanceResponse struct {
	Koinu Koinu
}

func (s *Service) GetBalance(ctx context.Context, req GetBalanceRequest) (*GetBalanceResponse, error) {
	s.balancesMutex.RLock()
	defer s.balancesMutex.RUnlock()

	koinu, ok := s.balances[req.Email]
	if !ok {
		return &GetBalanceResponse{}, nil
	}

	return &GetBalanceResponse{Koinu: koinu}, nil
}
