package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultTailscaleHostname = "reader"
	DefaultTailscaleStateDir = "tailscale"
	TailscaleAuthKeyEnvKey   = "TAILSCALE_AUTHKEY"
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
		t.AuthKey = os.Getenv(TailscaleAuthKeyEnvKey)
	}

	if t.StateDir == "" {
		t.StateDir = DefaultTailscaleStateDir
	}

	if !filepath.IsAbs(t.StateDir) {
		t.StateDir = filepath.Join(base, t.StateDir)
	}

	return nil
}
