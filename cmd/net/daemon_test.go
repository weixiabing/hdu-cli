package net

import "testing"

func TestNewDaemonCommandUse(t *testing.T) {
	cmd := newDaemonCmd()
	if cmd.Use != "daemon" {
		t.Fatalf("expected daemon command, got %q", cmd.Use)
	}
}
