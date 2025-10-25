package miniflux

import (
	"fmt"
	"net"

	"github.com/kellegous/poop"
	"miniflux.app/v2/client"
)

type Option func(*Server) error

func WithDatabase(
	name string,
	username string,
	password string,
) Option {
	return func(s *Server) error {
		s.env["DATABASE_URL"] = fmt.Sprintf(
			"user=%s password=%s dbname=%s sslmode=disable",
			username,
			password,
			name)
		return nil
	}
}

func toBaseURL(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	} else if port == "" {
		return "", poop.New("no port provided")
	}

	if host == "" {
		host = "localhost"
	}

	return fmt.Sprintf("http://%s:%s", host, port), nil
}

func WithListenAddress(addr string) Option {
	return func(s *Server) error {
		url, err := toBaseURL(addr)
		if err != nil {
			return poop.Chain(err)
		}
		s.env["LISTEN_ADDR"] = addr
		s.baseURL = url
		return nil
	}
}

func WithBaseURL(url string) Option {
	return func(s *Server) error {
		s.env["BASE_URL"] = url
		return nil
	}
}

func WithAdmin(username string, password string) Option {
	return func(s *Server) error {
		s.env["CREATE_ADMIN"] = "1"
		s.env["ADMIN_USERNAME"] = username
		s.env["ADMIN_PASSWORD"] = password

		s.opts = []client.Option{
			client.WithCredentials(username, password),
		}

		return nil
	}
}

func WithRunMigrations(v bool) Option {
	return func(s *Server) error {
		if v {
			s.env["RUN_MIGRATIONS"] = "1"
		} else {
			delete(s.env, "RUN_MIGRATIONS")
		}
		return nil
	}
}

func WithAuthProxyHeader(header string) Option {
	return func(s *Server) error {
		s.env["AUTH_PROXY_HEADER"] = header
		return nil
	}
}

func WithAuthProxyUserCreation(v bool) Option {
	return func(s *Server) error {
		s.env["AUTH_PROXY_USER_CREATION"] = "1"
		return nil
	}
}

func WithDebugLogging(v bool) Option {
	return func(s *Server) error {
		if v {
			s.env["DEBUG"] = "1"
		} else {
			delete(s.env, "DEBUG")
		}
		return nil
	}
}
