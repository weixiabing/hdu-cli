package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	appconfig "github.com/hduhelp/hdu-cli/internal/config"
	"github.com/hduhelp/hdu-cli/internal/core"
	"github.com/hduhelp/hdu-cli/pkg/srun"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type SettingsInput struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Endpoint      string `json:"endpoint"`
	ACID          string `json:"acid"`
	CheckInterval int    `json:"checkIntervalSeconds"`
	AutoConnect   bool   `json:"autoConnect"`
	AutoReconnect bool   `json:"autoReconnect"`
	LaunchAtLogin bool   `json:"launchAtLogin"`
}

type DesktopState struct {
	Phase         string `json:"phase"`
	Message       string `json:"message"`
	Online        bool   `json:"online"`
	Username      string `json:"username"`
	Endpoint      string `json:"endpoint"`
	ACID          string `json:"acid"`
	CheckInterval int    `json:"checkIntervalSeconds"`
	AutoConnect   bool   `json:"autoConnect"`
	AutoReconnect bool   `json:"autoReconnect"`
	LaunchAtLogin bool   `json:"launchAtLogin"`
}

type DesktopApp struct {
	mu      sync.RWMutex
	app     *application.App
	window  *application.WebviewWindow
	session *core.SessionService
	cfgPath string
	cfg     appconfig.AppConfig
	state   DesktopState
}

func NewDesktopApp(session *core.SessionService) *DesktopApp {
	cfg := appconfig.Default()
	return &DesktopApp{
		session: session,
		cfg:     cfg,
		state: DesktopState{
			Phase:         string(core.PhaseDisconnected),
			Endpoint:      cfg.Endpoint,
			ACID:          cfg.ACID,
			CheckInterval: cfg.CheckIntervalSeconds,
			AutoConnect:   cfg.AutoConnect,
			AutoReconnect: cfg.AutoReconnect,
		},
	}
}

func (d *DesktopApp) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	_ = ctx
	_ = options
	d.app = application.Get()

	path, err := defaultConfigPath()
	if err != nil {
		return err
	}
	d.cfgPath = path
	if err := d.reloadState(); err != nil {
		return err
	}
	return nil
}

func (d *DesktopApp) AttachWindow(window *application.WebviewWindow) {
	d.window = window
}

func (d *DesktopApp) CurrentState() DesktopState {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state
}

func (d *DesktopApp) ConnectNow() (DesktopState, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	creds := core.Credentials{
		Username: d.cfg.Username,
		Password: d.password(),
	}
	if creds.Username == "" || creds.Password == "" {
		d.state.Phase = string(core.PhaseFailed)
		d.state.Message = "username or password is empty"
		return d.state, errors.New(d.state.Message)
	}

	status, err := d.ensureSession().Login(context.Background(), creds)
	d.applyStatus(status, err)
	return d.state, err
}

func (d *DesktopApp) DisconnectNow() (DesktopState, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cfg.Username == "" {
		d.state.Phase = string(core.PhaseDisconnected)
		d.state.Online = false
		d.state.Message = ""
		return d.state, nil
	}

	err := d.ensureSession().Logout(context.Background(), core.Credentials{Username: d.cfg.Username})
	if err != nil {
		d.state.Phase = string(core.PhaseFailed)
		d.state.Message = err.Error()
		return d.state, err
	}
	d.state.Phase = string(core.PhaseDisconnected)
	d.state.Online = false
	d.state.Message = ""
	return d.state, nil
}

func (d *DesktopApp) SaveSettings(input SettingsInput) (DesktopState, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	cfg := appconfig.Default()
	if input.Endpoint != "" {
		cfg.Endpoint = input.Endpoint
	}
	if input.ACID != "" {
		cfg.ACID = input.ACID
	}
	cfg.Username = input.Username
	if input.CheckInterval > 0 {
		cfg.CheckIntervalSeconds = input.CheckInterval
	}
	cfg.AutoConnect = input.AutoConnect
	cfg.AutoReconnect = input.AutoReconnect

	if d.cfgPath == "" {
		path, err := defaultConfigPath()
		if err != nil {
			return d.state, err
		}
		d.cfgPath = path
	}
	if err := appconfig.Save(d.cfgPath, cfg, input.Password); err != nil {
		return d.state, err
	}
	d.cfg = cfg
	d.state.Username = cfg.Username
	d.state.Endpoint = cfg.Endpoint
	d.state.ACID = cfg.ACID
	d.state.CheckInterval = cfg.CheckIntervalSeconds
	d.state.AutoConnect = cfg.AutoConnect
	d.state.AutoReconnect = cfg.AutoReconnect
	d.state.LaunchAtLogin = input.LaunchAtLogin
	d.state.Message = "settings saved"

	if err := d.applyLaunchAtLogin(input.LaunchAtLogin); err != nil {
		d.state.Message = err.Error()
		return d.state, err
	}
	return d.state, nil
}

func (d *DesktopApp) ToggleLaunchAtLogin(enabled bool) (DesktopState, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.applyLaunchAtLogin(enabled); err != nil {
		d.state.Message = err.Error()
		return d.state, err
	}
	d.state.LaunchAtLogin = enabled
	d.state.Message = "launch at login updated"
	return d.state, nil
}

func (d *DesktopApp) ShowWindow() {
	if d.window != nil {
		d.window.Show()
	}
}

func (d *DesktopApp) HideWindow() {
	if d.window != nil {
		d.window.Hide()
	}
}

func (d *DesktopApp) ReconnectStatus() (DesktopState, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	status, err := d.ensureSession().CurrentStatus(context.Background())
	d.applyStatus(status, err)
	return d.state, err
}

func (d *DesktopApp) ensureSession() *core.SessionService {
	if d.session != nil {
		return d.session
	}
	d.session = core.NewSessionService(srun.New(d.cfg.Endpoint, d.cfg.ACID))
	return d.session
}

func (d *DesktopApp) applyStatus(status core.Status, err error) {
	if status.Phase == "" {
		status.Phase = core.PhaseDisconnected
	}
	d.state.Phase = string(status.Phase)
	d.state.Online = status.Online
	if status.Message != "" {
		d.state.Message = status.Message
	} else if err != nil {
		d.state.Message = err.Error()
	} else {
		d.state.Message = ""
	}
}

func (d *DesktopApp) reloadState() error {
	cfg, err := appconfig.Load(d.cfgPath)
	if err != nil {
		return err
	}
	d.cfg = cfg
	d.state.Username = cfg.Username
	d.state.Endpoint = cfg.Endpoint
	d.state.ACID = cfg.ACID
	d.state.CheckInterval = cfg.CheckIntervalSeconds
	d.state.AutoConnect = cfg.AutoConnect
	d.state.AutoReconnect = cfg.AutoReconnect
	d.session = core.NewSessionService(srun.New(cfg.Endpoint, cfg.ACID))

	if d.app != nil {
		st, statusErr := d.app.Autostart.Status()
		if statusErr == nil {
			d.state.LaunchAtLogin = st.Enabled
		}
	}

	if cfg.AutoConnect && cfg.Username != "" && d.password() != "" {
		go func() {
			time.Sleep(500 * time.Millisecond)
			_, _ = d.ConnectNow()
		}()
	}

	return nil
}

func (d *DesktopApp) password() string {
	if d.cfgPath == "" {
		return ""
	}
	v := viper.New()
	v.SetConfigFile(d.cfgPath)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return ""
	}
	if password := v.GetString("password"); password != "" {
		return password
	}
	if password := v.GetString("net.auth.password"); password != "" {
		return password
	}
	return ""
}

func defaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, ".hdu-cli.yaml"), nil
}

func (d *DesktopApp) applyLaunchAtLogin(enabled bool) error {
	if d.app == nil {
		return nil
	}
	if enabled {
		return d.app.Autostart.Enable()
	}
	return d.app.Autostart.Disable()
}
