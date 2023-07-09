package config

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/kellegous/reader/pkg/logging"
	"go.uber.org/zap"
)

const (
	DefaultDataDir  = "db"
	DefaultUsername = "reader"
)

type Postgres struct {
	DataDir  string `yaml:"data-dir"`
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

	if p.Password == "" {
		return errors.New("postgres.password is required")
	}

	if !filepath.IsAbs(p.DataDir) {
		logging.L(context.Background()).Info("making data dir absolute",
			zap.String("base", base))
		p.DataDir = filepath.Join(base, p.DataDir)
	}

	return nil
}
