package config

import "errors"

type Miniflux struct {
	AdminUsername string `yaml:"admin-username"`
	AdminPassword string `yaml:"admin-password"`
	AutoLoginAs   string `yaml:"auto-login-as"`
}

func (m *Miniflux) apply() error {
	if m.AdminUsername == "" {
		return errors.New("miniflux.admin-username is required")
	}
	if m.AdminPassword == "" {
		return errors.New("miniflux.admin-password is required")
	}
	return nil
}
