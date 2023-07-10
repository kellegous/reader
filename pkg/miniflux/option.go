package miniflux

import "fmt"

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

func WithListenAddress(addr string) Option {
	return func(s *Server) error {
		s.env["LISTEN_ADDR"] = addr
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
