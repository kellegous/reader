package miniflux

import (
	"context"
	"os"
	"os/exec"

	"miniflux.app/v2/client"
)

type Server struct {
	proc *os.Process
	opts Options
}

func (s *Server) Stop() error {
	return s.proc.Kill()
}

func (s *Server) Client(opts ...client.Option) *client.Client {
	if a := s.opts.admin; a != nil {
		opts = append(opts, client.WithCredentials(a.username, a.password))
	}

	return client.NewClientWithOptions(
		s.opts.externalURL,
		opts...)
}

func (s *Server) BaseURL() string {
	return s.opts.internalURL
}

func Start(ctx context.Context, opts ...Option) (*Server, error) {
	s := &Server{}
	for _, opt := range opts {
		if err := opt(&s.opts); err != nil {
			return nil, err
		}
	}

	c := exec.CommandContext(ctx, "miniflux")
	c.Env = s.opts.env()
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		return nil, err
	}

	s.proc = c.Process

	return s, nil
}
