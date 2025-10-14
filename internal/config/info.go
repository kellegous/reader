package config

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Info struct {
	Postgres Postgres `yaml:"postgres"`
	Miniflux Miniflux `yaml:"miniflux"`
	Web      Web      `yaml:"web"`
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

	if err := n.Postgres.apply(base); err != nil {
		return err
	}

	if err := n.Miniflux.apply(); err != nil {
		return err
	}

	if err := n.Web.apply(); err != nil {
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
