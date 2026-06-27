package core

import (
	"context"
	"errors"
)

type Phase string

const (
	PhaseDisconnected Phase = "disconnected"
	PhaseConnecting   Phase = "connecting"
	PhaseConnected    Phase = "connected"
	PhaseRetrying     Phase = "retrying"
	PhaseFailed       Phase = "failed"
)

var ErrAuthenticationFailed = errors.New("authentication failed")

type Credentials struct {
	Username string
	Password string
}

type Status struct {
	Phase   Phase
	Online  bool
	Message string
}

type PortalClient interface {
	Login(ctx context.Context, username, password string) (Status, error)
	Logout(ctx context.Context, username string) error
	CurrentStatus(ctx context.Context) (Status, error)
	InternetReachable(ctx context.Context) bool
}
