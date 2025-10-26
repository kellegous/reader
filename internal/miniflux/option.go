package miniflux

import (
	"fmt"
	"net"

	"github.com/kellegous/poop"
)

type Option func(*Options) error

type Options struct {
	databaseURL           string
	externalURL           string
	internalURL           string
	listenAddr            string
	admin                 *login
	runMigrations         bool
	authProxyHeader       string
	authProxyUserCreation bool
	debugLogging          bool
}

type login struct {
	username string
	password string
}

func (o *Options) env() []string {
	var env []string
	if o.databaseURL != "" {
		env = append(env, "DATABASE_URL="+o.databaseURL)
	}
	if o.listenAddr != "" {
		env = append(env, "LISTEN_ADDR="+o.listenAddr)
	}
	if o.externalURL != "" {
		env = append(env, "BASE_URL="+o.externalURL)
	}
	if o.listenAddr != "" {
		env = append(env, "LISTEN_ADDR="+o.listenAddr)
	}
	if a := o.admin; a != nil {
		env = append(env, "CREATE_ADMIN=1")
		env = append(env, "ADMIN_USERNAME="+a.username)
		env = append(env, "ADMIN_PASSWORD="+a.password)
	}
	if o.runMigrations {
		env = append(env, "RUN_MIGRATIONS=1")
	}
	if o.authProxyHeader != "" {
		env = append(env, "AUTH_PROXY_HEADER="+o.authProxyHeader)
	}
	if o.authProxyUserCreation {
		env = append(env, "AUTH_PROXY_USER_CREATION=1")
	}
	if o.debugLogging {
		env = append(env, "DEBUG=1")
	}
	return env
}

func WithDatabase(
	name string,
	username string,
	password string,
) Option {
	return func(o *Options) error {
		o.databaseURL = fmt.Sprintf(
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
	return func(o *Options) error {
		url, err := toBaseURL(addr)
		if err != nil {
			return poop.Chain(err)
		}
		o.listenAddr = addr
		o.internalURL = url
		return nil
	}
}

func WithBaseURL(url string) Option {
	return func(o *Options) error {
		o.externalURL = url
		return nil
	}
}

func WithAdmin(username string, password string) Option {
	return func(o *Options) error {
		o.admin = &login{
			username: username,
			password: password,
		}
		return nil
	}
}

func WithRunMigrations(v bool) Option {
	return func(o *Options) error {
		o.runMigrations = v
		return nil
	}
}

func WithAuthProxyHeader(header string) Option {
	return func(o *Options) error {
		o.authProxyHeader = header
		return nil
	}
}

func WithAuthProxyUserCreation(v bool) Option {
	return func(o *Options) error {
		o.authProxyUserCreation = v
		return nil
	}
}

func WithDebugLogging(v bool) Option {
	return func(o *Options) error {
		o.debugLogging = v
		return nil
	}
}
