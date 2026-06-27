package core

import (
	"context"
	"time"
)

type ReconnectConfig struct {
	Interval   time.Duration
	MaxBackoff time.Duration
}

type ReconnectManager struct {
	session *SessionService
	cfg     ReconnectConfig
	sleep   func(time.Duration)
}

func NewReconnectManager(session *SessionService, cfg ReconnectConfig) *ReconnectManager {
	return &ReconnectManager{
		session: session,
		cfg:     cfg,
		sleep:   time.Sleep,
	}
}

func (m *ReconnectManager) ReconnectOnce(ctx context.Context, creds Credentials) Status {
	backoff := m.cfg.Interval
	if backoff <= 0 {
		backoff = time.Millisecond
	}

	for {
		status, err := m.session.Login(ctx, creds)
		if err == nil {
			return status
		}

		if m.cfg.MaxBackoff > 0 && backoff < m.cfg.MaxBackoff {
			backoff *= 2
			if backoff > m.cfg.MaxBackoff {
				backoff = m.cfg.MaxBackoff
			}
		}

		m.sleep(backoff)
	}
}
