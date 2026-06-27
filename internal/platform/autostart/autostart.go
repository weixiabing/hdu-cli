package autostart

type Manager interface {
	Enable() error
	Disable() error
	Enabled() (bool, error)
}
