package core

import "context"

type SessionService struct {
	client PortalClient
}

func NewSessionService(client PortalClient) *SessionService {
	return &SessionService{client: client}
}

func (s *SessionService) Login(ctx context.Context, creds Credentials) (Status, error) {
	return s.client.Login(ctx, creds.Username, creds.Password)
}
