package core

import (
	"context"
	"errors"
	"time"
)

type ReconnectConfig struct {
	Interval   time.Duration
	MaxBackoff time.Duration
}

type ReconnectManager struct {
	session      *SessionService
	cfg          ReconnectConfig
	sleep        SleepFunc
	stateHandler StateHandler
}

func NewReconnectManager(session *SessionService, cfg ReconnectConfig) *ReconnectManager {
	return &ReconnectManager{
		session: session,
		cfg:     cfg,
		sleep:   waitForBackoff,
	}
}

func (m *ReconnectManager) OnStateChange(handler StateHandler) {
	m.stateHandler = handler
}

func (m *ReconnectManager) SetSleep(sleep SleepFunc) {
	m.sleep = sleep
}

func (m *ReconnectManager) ReconnectOnce(ctx context.Context, creds Credentials) (Status, error) {
	backoff := m.cfg.Interval
	if backoff <= 0 {
		backoff = time.Millisecond
	}

	for {
		if err := ctx.Err(); err != nil {
			status := Status{Phase: PhaseFailed, Online: false, Message: err.Error()}
			m.emit(status)
			return status, err
		}

		status, err := m.session.Login(ctx, creds)
		if err == nil {
			if status.Phase == "" {
				status.Phase = PhaseConnected
			}
			m.emit(status)
			return status, nil
		}

		if errors.Is(err, ErrAuthenticationFailed) {
			failed := Status{Phase: PhaseFailed, Online: false, Message: status.Message}
			if failed.Message == "" {
				failed.Message = err.Error()
			}
			m.emit(failed)
			return failed, err
		}

		retrying := Status{Phase: PhaseRetrying, Online: false, Message: status.Message}
		if retrying.Message == "" {
			retrying.Message = err.Error()
		}
		m.emit(retrying)

		if err := m.sleep(ctx, backoff); err != nil {
			failed := Status{Phase: PhaseFailed, Online: false, Message: err.Error()}
			m.emit(failed)
			return failed, err
		}

		if m.cfg.MaxBackoff > 0 && backoff < m.cfg.MaxBackoff {
			backoff *= 2
			if backoff > m.cfg.MaxBackoff {
				backoff = m.cfg.MaxBackoff
			}
		}

	}
}

func (m *ReconnectManager) emit(status Status) {
	if m.stateHandler != nil {
		m.stateHandler(status)
	}
}

func waitForBackoff(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
