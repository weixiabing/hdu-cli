package autostart

type noopManager struct{}

func NewManager() Manager {
	return &noopManager{}
}

func (n *noopManager) Enable() error {
	return nil
}

func (n *noopManager) Disable() error {
	return nil
}

func (n *noopManager) Enabled() (bool, error) {
	return false, nil
}
