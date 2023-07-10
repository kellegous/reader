package miniflux

import (
	"context"
	"os"
	"os/exec"
)

type Server struct {
	env  map[string]string
	proc *os.Process
}

func (s *Server) Stop() error {
	return s.proc.Kill()
}

func Start(ctx context.Context, opts ...Option) (*Server, error) {
	s := &Server{
		env: map[string]string{},
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	c := exec.CommandContext(ctx, "miniflux")
	c.Env = s.getEnv()
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		return nil, err
	}

	s.proc = c.Process

	return s, nil
}

func (s *Server) getEnv() []string {
	env := make([]string, 0, len(s.env))
	for k, v := range s.env {
		env = append(env, k+"="+v)
	}
	return env
}
