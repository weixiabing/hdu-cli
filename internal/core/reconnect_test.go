package core

import (
	"context"
	"testing"
	"time"
)

type fakeReconnectClient struct {
	loginFn func(ctx context.Context, username, password string) (Status, error)
}

func (f *fakeReconnectClient) Login(ctx context.Context, username, password string) (Status, error) {
	return f.loginFn(ctx, username, password)
}

func (f *fakeReconnectClient) Logout(ctx context.Context, username string) error { return nil }

func (f *fakeReconnectClient) CurrentStatus(ctx context.Context) (Status, error) {
	return Status{Phase: PhaseDisconnected, Online: false}, nil
}

func (f *fakeReconnectClient) InternetReachable(ctx context.Context) bool { return true }

func TestReconnectManagerRetriesUntilConnected(t *testing.T) {
	attempts := 0
	client := &fakeReconnectClient{
		loginFn: func(ctx context.Context, username, password string) (Status, error) {
			attempts++
			if attempts < 3 {
				return Status{Phase: PhaseFailed, Online: false}, ErrAuthenticationFailed
			}
			return Status{Phase: PhaseConnected, Online: true}, nil
		},
	}

	manager := NewReconnectManager(NewSessionService(client), ReconnectConfig{
		Interval:   time.Millisecond,
		MaxBackoff: 5 * time.Millisecond,
	})

	status := manager.ReconnectOnce(context.Background(), Credentials{
		Username: "20230001",
		Password: "secret",
	})

	if status.Phase != PhaseConnected || attempts != 3 {
		t.Fatalf("expected third attempt success, got %#v with %d attempts", status, attempts)
	}
}
