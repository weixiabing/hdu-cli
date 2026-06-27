package service

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

type UnitConfig struct {
	BinaryPath string
	ConfigPath string
}

const unitName = "hdu-cli.service"

const userUnit = `[Unit]
Description=HDU campus network client
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart={{ .BinaryPath }} net daemon --config {{ .ConfigPath }}
Restart=always
RestartSec=5

[Install]
WantedBy=default.target
`

func RenderUserUnit(cfg UnitConfig) (string, error) {
	var out bytes.Buffer
	tpl, err := template.New("systemd-user-unit").Parse(userUnit)
	if err != nil {
		return "", err
	}
	if err := tpl.Execute(&out, cfg); err != nil {
		return "", err
	}
	return out.String(), nil
}

func UserUnitPath(home string) string {
	return filepath.Join(home, ".config", "systemd", "user", unitName)
}

func InstallUserUnit(home string, cfg UnitConfig) (string, error) {
	unit, err := RenderUserUnit(cfg)
	if err != nil {
		return "", err
	}

	path := UserUnitPath(home)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(unit), 0o600); err != nil {
		return "", err
	}
	return path, nil
}
