package service

import (
	"strings"
	"testing"
)

func TestRenderUserUnit(t *testing.T) {
	unit, err := RenderUserUnit(UnitConfig{
		BinaryPath: "/home/test/.local/bin/hdu-cli",
		ConfigPath: "/home/test/.hdu-cli.yaml",
	})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if !strings.Contains(unit, "ExecStart=/home/test/.local/bin/hdu-cli net daemon --config /home/test/.hdu-cli.yaml") {
		t.Fatalf("missing daemon ExecStart: %s", unit)
	}
}
