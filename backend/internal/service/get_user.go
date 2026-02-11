package service

import "context"

func (s *Service) GetUser(ctx context.Context, authToken AuthToken) (*User, error) {
	s.sessionsMutex.RLock()
	defer s.sessionsMutex.RUnlock()

	user, ok := s.sessions[authToken]
	if !ok {
		return nil, nil
	}

	return user, nil
}
