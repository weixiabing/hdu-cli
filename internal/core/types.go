package core

import "errors"

type Phase string

const (
	PhaseDisconnected Phase = "disconnected"
	PhaseConnecting   Phase = "connecting"
	PhaseConnected    Phase = "connected"
	PhaseRetrying     Phase = "retrying"
	PhaseFailed       Phase = "failed"
)

var ErrAuthenticationFailed = errors.New("authentication failed")
var ErrReconnectCanceled = errors.New("reconnect canceled")

type Credentials struct {
	Username string
	Password string
}

type Status struct {
	Phase   Phase
	Online  bool
	Message string
}
