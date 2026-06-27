package core

import (
	"context"
	"testing"
)

type fakePortalClient struct {
	loginCalls int
}

func (f *fakePortalClient) Login(ctx context.Context, username, password string) (Status, error) {
	f.loginCalls++
	return Status{Phase: PhaseConnected, Online: true}, nil
}

func (f *fakePortalClient) Logout(ctx context.Context, username string) error { return nil }

func (f *fakePortalClient) CurrentStatus(ctx context.Context) (Status, error) {
	return Status{Phase: PhaseDisconnected, Online: false}, nil
}

func (f *fakePortalClient) InternetReachable(ctx context.Context) bool { return true }

func TestSessionServiceLogin(t *testing.T) {
	client := &fakePortalClient{}
	service := NewSessionService(client)

	status, err := service.Login(context.Background(), Credentials{
		Username: "20230001",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("login returned error: %v", err)
	}
	if !status.Online || status.Phase != PhaseConnected {
		t.Fatalf("expected connected status, got %#v", status)
	}
	if client.loginCalls != 1 {
		t.Fatalf("expected one login call, got %d", client.loginCalls)
	}
}
