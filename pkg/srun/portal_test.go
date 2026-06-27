package srun

import (
	"context"
	"errors"
	"testing"

	"github.com/hduhelp/hdu-cli/internal/core"
)

func TestPortalServerLoginFetchesChallengeBeforePortalLogin(t *testing.T) {
	var calls []string
	server := &PortalServer{
		userInfo: &userInfo{OnlineIp: "10.0.0.1"},
		adapterHooks: portalAdapterHooks{
			getUserInfo: func() (*userInfo, error) {
				calls = append(calls, "userInfo")
				return &userInfo{OnlineIp: "10.0.0.1", Error: "ok"}, nil
			},
			getChallenge: func() (*challenge, error) {
				calls = append(calls, "challenge")
				return &challenge{Challenge: "token", Error: "ok"}, nil
			},
			portalLogin: func() (*loginResponse, error) {
				calls = append(calls, "login")
				return &loginResponse{Error: "ok"}, nil
			},
		},
	}

	status, err := server.Login(context.Background(), "20230001", "secret")
	if err != nil {
		t.Fatalf("login returned error: %v", err)
	}
	if status.Phase != core.PhaseConnected || !status.Online {
		t.Fatalf("expected connected status, got %#v", status)
	}
	if len(calls) != 3 || calls[0] != "userInfo" || calls[1] != "challenge" || calls[2] != "login" {
		t.Fatalf("expected userInfo -> challenge -> login order, got %#v", calls)
	}
}

func TestPortalServerCurrentStatusReturnsDisconnectedWithoutErrorWhenOffline(t *testing.T) {
	server := &PortalServer{
		adapterHooks: portalAdapterHooks{
			getUserInfo: func() (*userInfo, error) {
				return &userInfo{Error: "not_online_error"}, nil
			},
		},
	}

	status, err := server.CurrentStatus(context.Background())
	if err != nil {
		t.Fatalf("expected offline status without error, got %v", err)
	}
	if status.Phase != core.PhaseDisconnected || status.Online {
		t.Fatalf("expected disconnected offline status, got %#v", status)
	}
}

func TestPortalServerLoginMapsAuthenticationFailure(t *testing.T) {
	server := &PortalServer{
		userInfo: &userInfo{OnlineIp: "10.0.0.1"},
		adapterHooks: portalAdapterHooks{
			getUserInfo: func() (*userInfo, error) {
				return &userInfo{OnlineIp: "10.0.0.1", Error: "ok"}, nil
			},
			getChallenge: func() (*challenge, error) {
				return &challenge{Challenge: "token", Error: "ok"}, nil
			},
			portalLogin: func() (*loginResponse, error) {
				return &loginResponse{Error: "fail", ErrorMsg: "password is incorrect"}, nil
			},
		},
	}

	status, err := server.Login(context.Background(), "20230001", "secret")
	if !errors.Is(err, core.ErrAuthenticationFailed) {
		t.Fatalf("expected authentication failure, got %v", err)
	}
	if status.Phase != core.PhaseFailed {
		t.Fatalf("expected failed status, got %#v", status)
	}
}

func TestPortalServerLogoutReturnsNilWhenAlreadyOffline(t *testing.T) {
	logoutCalled := false
	server := &PortalServer{
		adapterHooks: portalAdapterHooks{
			getUserInfo: func() (*userInfo, error) {
				return &userInfo{Error: "not_online_error"}, nil
			},
			portalLogout: func() (*logoutResponse, error) {
				logoutCalled = true
				return &logoutResponse{Error: "ok"}, nil
			},
		},
	}

	if err := server.Logout(context.Background(), "20230001"); err != nil {
		t.Fatalf("expected offline logout to return nil, got %v", err)
	}
	if logoutCalled {
		t.Fatal("expected logout request to be skipped when already offline")
	}
}
