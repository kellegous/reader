package miniflux

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kellegous/poop"
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

func (s *Server) SQLConn(ctx context.Context) (*pgx.Conn, error) {
	c, err := pgx.Connect(ctx, s.opts.databaseURL)
	if err != nil {
		return nil, poop.Chain(err)
	}
	return c, nil
}

func (s *Server) BaseURL() string {
	return s.opts.internalURL
}

func (s *Server) provisionAuthProxyUser(
	ctx context.Context,
	user string,
) error {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return poop.Chain(err)
	}

	// this check has to have a cookie jar so that miniflux can see
	// its own session cookie, otherwise it will redirect indefinitely
	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.BaseURL()+"/", nil)
	if err != nil {
		return poop.Chain(err)
	}

	req.Header.Set(s.opts.authProxy.header, user)

	res, err := client.Do(req)
	if err != nil {
		return poop.Chain(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return poop.Newf("status %d for auth proxy user: %s", res.StatusCode, user)
	}

	return nil
}

func waitForLiveness(
	ctx context.Context,
	s *Server,
	timeout time.Duration,
) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	doCheck := func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.BaseURL()+"/liveness", nil)
		if err != nil {
			return poop.Chain(err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return poop.Chain(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return poop.Newf("status %d for liveness check", res.StatusCode)
		}

		return nil
	}

	for {
		if err := doCheck(); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
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

	if err := waitForLiveness(ctx, s, 10*time.Second); err != nil {
		return nil, poop.Chain(err)
	}

	if p := s.opts.authProxy; p != nil {
		for _, user := range p.users {
			if err := s.provisionAuthProxyUser(ctx, user); err != nil {
				return nil, poop.Chain(err)
			}
		}
	}

	return s, nil
}
