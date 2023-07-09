package config

import "errors"

const (
	DefaultTailScaleHostname = "reader"
)

type Tailscale struct {
	AuthKey  string `yaml:"auth-key"`
	Hostname string `yaml:"hostname"`
}

func (t *Tailscale) apply() error {
	if t.Hostname == "" {
		t.Hostname = DefaultTailScaleHostname
	}

	if t.AuthKey == "" {
		return errors.New("tailscale.auth-key is required")
	}

	return nil
}
