package config

import (
	"errors"
	"path/filepath"
)

const (
	DefaultTailscaleHostname = "reader"
	DefaultTailscaleStateDir = "tailscale"
)

type Tailscale struct {
	AuthKey  string `yaml:"auth-key"`
	Hostname string `yaml:"hostname"`
	StateDir string `yaml:"state-dir"`
}

func (t *Tailscale) apply(base string) error {
	if t.Hostname == "" {
		t.Hostname = DefaultTailscaleHostname
	}

	if t.AuthKey == "" {
		return errors.New("tailscale.auth-key is required")
	}

	if t.StateDir == "" {
		t.StateDir = DefaultTailscaleStateDir
	}

	if !filepath.IsAbs(t.StateDir) {
		t.StateDir = filepath.Join(base, t.StateDir)
	}

	return nil
}
