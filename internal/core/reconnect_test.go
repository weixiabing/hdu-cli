package core

import (
	"context"
	"errors"
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
				return Status{Phase: PhaseRetrying, Online: false, Message: "temporary failure"}, errors.New("temporary failure")
			}
			return Status{Phase: PhaseConnected, Online: true}, nil
		},
	}

	manager := NewReconnectManager(NewSessionService(client), ReconnectConfig{
		Interval:   time.Millisecond,
		MaxBackoff: 5 * time.Millisecond,
	})

	var observed []Phase
	manager.OnStateChange(func(status Status) {
		observed = append(observed, status.Phase)
	})

	status, err := manager.ReconnectOnce(context.Background(), Credentials{
		Username: "20230001",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("reconnect returned error: %v", err)
	}

	if status.Phase != PhaseConnected || attempts != 3 {
		t.Fatalf("expected third attempt success, got %#v with %d attempts", status, attempts)
	}
	if len(observed) < 3 {
		t.Fatalf("expected state transitions, got %#v", observed)
	}
	if observed[0] != PhaseRetrying {
		t.Fatalf("expected first state retrying, got %#v", observed)
	}
	if observed[len(observed)-1] != PhaseConnected {
		t.Fatalf("expected final state connected, got %#v", observed)
	}
}

func TestReconnectManagerAuthenticationFailureStopsRetrying(t *testing.T) {
	attempts := 0
	client := &fakeReconnectClient{
		loginFn: func(ctx context.Context, username, password string) (Status, error) {
			attempts++
			return Status{Phase: PhaseFailed, Online: false, Message: "bad credentials"}, ErrAuthenticationFailed
		},
	}

	manager := NewReconnectManager(NewSessionService(client), ReconnectConfig{
		Interval:   time.Millisecond,
		MaxBackoff: 5 * time.Millisecond,
	})

	status, err := manager.ReconnectOnce(context.Background(), Credentials{
		Username: "20230001",
		Password: "secret",
	})
	if !errors.Is(err, ErrAuthenticationFailed) {
		t.Fatalf("expected authentication error, got %v", err)
	}
	if attempts != 1 {
		t.Fatalf("expected one attempt, got %d", attempts)
	}
	if status.Phase != PhaseFailed {
		t.Fatalf("expected failed status, got %#v", status)
	}
}

func TestReconnectManagerStopsOnContextCancellation(t *testing.T) {
	attempts := 0
	client := &fakeReconnectClient{
		loginFn: func(ctx context.Context, username, password string) (Status, error) {
			attempts++
			return Status{Phase: PhaseRetrying, Online: false, Message: "temporary failure"}, errors.New("temporary failure")
		},
	}

	manager := NewReconnectManager(NewSessionService(client), ReconnectConfig{
		Interval:   time.Millisecond,
		MaxBackoff: 5 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())
	manager.SetSleep(func(_ context.Context, _ time.Duration) error {
		cancel()
		return ctx.Err()
	})

	status, err := manager.ReconnectOnce(ctx, Credentials{
		Username: "20230001",
		Password: "secret",
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled error, got %v", err)
	}
	if attempts != 1 {
		t.Fatalf("expected one attempt before cancellation, got %d", attempts)
	}
	if status.Phase != PhaseFailed {
		t.Fatalf("expected failed status on cancellation, got %#v", status)
	}
}
