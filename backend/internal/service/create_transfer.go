package service

import (
	"context"
	"fmt"
)

type CreateTransferRequest struct {
	SenderEmail    Email
	RecipientEmail Email
	AmountKoinu    Koinu
}

type CreateTransferResponse struct{}

func (s *Service) CreateTransfer(ctx context.Context, req CreateTransferRequest) (*CreateTransferResponse, error) {
	s.balancesMutex.Lock()
	defer s.balancesMutex.Unlock()

	senderBalance := s.balances[req.SenderEmail]

	if senderBalance < req.AmountKoinu {
		return nil, fmt.Errorf("%w: sender balance is less than amount", ErrInsufficientBalance)
	}

	s.balances[req.SenderEmail] -= req.AmountKoinu
	s.balances[req.RecipientEmail] += req.AmountKoinu

	return &CreateTransferResponse{}, nil
}
