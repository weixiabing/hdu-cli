package main

import "testing"

func TestNewDesktopAppHasDefaultState(t *testing.T) {
	app := NewDesktopApp(nil)
	state := app.CurrentState()
	if state.Phase != "disconnected" {
		t.Fatalf("expected disconnected phase, got %q", state.Phase)
	}
}
