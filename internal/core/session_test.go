package core

import (
	"context"
	"errors"
	"testing"
)

type fakePortalClient struct {
	loginCalls         int
	logoutCalls        int
	currentStatusCalls int
	logoutErr          error
	currentStatus      Status
	currentStatusErr   error
}

func (f *fakePortalClient) Login(ctx context.Context, username, password string) (Status, error) {
	f.loginCalls++
	return Status{Phase: PhaseConnected, Online: true}, nil
}

func (f *fakePortalClient) Logout(ctx context.Context, username string) error {
	f.logoutCalls++
	return f.logoutErr
}

func (f *fakePortalClient) CurrentStatus(ctx context.Context) (Status, error) {
	f.currentStatusCalls++
	if f.currentStatus == (Status{}) && f.currentStatusErr == nil {
		return Status{Phase: PhaseDisconnected, Online: false}, nil
	}
	return f.currentStatus, f.currentStatusErr
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

func TestSessionServiceLogout(t *testing.T) {
	client := &fakePortalClient{}
	service := NewSessionService(client)

	if err := service.Logout(context.Background(), Credentials{Username: "20230001"}); err != nil {
		t.Fatalf("logout returned error: %v", err)
	}
	if client.logoutCalls != 1 {
		t.Fatalf("expected one logout call, got %d", client.logoutCalls)
	}
}

func TestSessionServiceCurrentStatus(t *testing.T) {
	client := &fakePortalClient{
		currentStatus: Status{Phase: PhaseConnected, Online: true},
	}
	service := NewSessionService(client)

	status, err := service.CurrentStatus(context.Background())
	if err != nil {
		t.Fatalf("current status returned error: %v", err)
	}
	if status.Phase != PhaseConnected || !status.Online {
		t.Fatalf("expected connected status, got %#v", status)
	}
	if client.currentStatusCalls != 1 {
		t.Fatalf("expected one current status call, got %d", client.currentStatusCalls)
	}
}

func TestSessionServiceCurrentStatusReturnsClientError(t *testing.T) {
	client := &fakePortalClient{
		currentStatusErr: errors.New("status unavailable"),
	}
	service := NewSessionService(client)

	_, err := service.CurrentStatus(context.Background())
	if err == nil {
		t.Fatal("expected current status error")
	}
}
