package config

import "errors"

type Miniflux struct {
	AdminUsername string `yaml:"admin-username"`
	AdminPassword string `yaml:"admin-password"`
}

func (m *Miniflux) apply(base string) error {
	if m.AdminUsername == "" {
		return errors.New("miniflux.admin-username is required")
	}
	if m.AdminPassword == "" {
		return errors.New("miniflux.admin-password is required")
	}
	return nil
}
