package core

import (
	"context"
	"time"
)

// PortalClient is the shared boundary that platform-specific portal adapters
// implement so the core session and reconnect services can stay transport-agnostic.
type PortalClient interface {
	Login(ctx context.Context, username, password string) (Status, error)
	Logout(ctx context.Context, username string) error
	CurrentStatus(ctx context.Context) (Status, error)
	InternetReachable(ctx context.Context) bool
}

type SleepFunc func(ctx context.Context, delay time.Duration) error

type StateHandler func(Status)
