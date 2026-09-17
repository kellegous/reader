package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/kellegous/glue/fn"
	"gopkg.in/yaml.v3"
)

type Info struct {
	Postgres Postgres `yaml:"postgres"`
	Miniflux Miniflux `yaml:"miniflux"`
	Web      Web      `yaml:"web"`
	Ollama   Ollama   `yaml:"ollama"`
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

	n.Ollama.apply()

	return nil
}

func (n *Info) ReadFile(src string) (err error) {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fn.WithCare(r.Close, &err)
	return n.Read(r, filepath.Dir(src))
}
