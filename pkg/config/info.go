package config

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Info struct {
	Tailscale Tailscale `yaml:"tailscale"`
	Postgres  Postgres  `yaml:"postgres"`
}

func (n *Info) Read(r io.Reader, base string) error {
	var err error
	if !filepath.IsAbs(base) {
		base, err = filepath.Abs(base)
		if err != nil {
			return err
		}
	}

	if err := yaml.NewDecoder(r).Decode(n); err != nil {
		return err
	}

	if err := n.Tailscale.apply(); err != nil {
		return err
	}

	if err := n.Postgres.apply(base); err != nil {
		return err
	}

	return nil
}

func (n *Info) ReadFile(src string) error {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()
	return n.Read(r, filepath.Dir(src))
}
