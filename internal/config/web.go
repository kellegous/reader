package config

import "errors"

const (
	DefaultWebAddr = ":4040"
)

type Web struct {
	Addr     string `yaml:"addr"`
	Hostname string `yaml:"hostname"`
}

func (w *Web) apply() error {
	if w.Addr == "" {
		w.Addr = DefaultWebAddr
	}

	if w.Hostname == "" {
		return errors.New("web.hostname is required")
	}

	return nil
}
