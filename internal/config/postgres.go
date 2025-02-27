package config

import (
	"errors"
	"path/filepath"
)

const (
	DefaultDataDir  = "db"
	DefaultUsername = "reader"
	DefaultDatabase = "reader"
)

type Postgres struct {
	DataDir  string `yaml:"data-dir"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (p *Postgres) apply(base string) error {
	if p.DataDir == "" {
		p.DataDir = DefaultDataDir
	}

	if p.Username == "" {
		p.Username = DefaultUsername
	}

	if p.Database == "" {
		p.Database = DefaultDatabase
	}

	if p.Password == "" {
		return errors.New("postgres.password is required")
	}

	if !filepath.IsAbs(p.DataDir) {
		p.DataDir = filepath.Join(base, p.DataDir)
	}

	return nil
}
