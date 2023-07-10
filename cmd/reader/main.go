package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/kellegous/tsweb"
	"go.uber.org/zap"
	"tailscale.com/ipn/ipnstate"
	"tailscale.com/tsnet"

	"github.com/kellegous/reader/pkg/config"
	"github.com/kellegous/reader/pkg/logging"
	"github.com/kellegous/reader/pkg/miniflux"
	"github.com/kellegous/reader/pkg/postgres"
	"github.com/kellegous/reader/pkg/web"
)

const (
	backendAddr = "127.0.0.1:9090"
)

type Flags struct {
	ConfigFile string
}

func (f *Flags) Register(fs *flag.FlagSet) {
	fs.StringVar(
		&f.ConfigFile,
		"config-file",
		"reader.yaml",
		"Path to the config file",
	)
}

func ensurePostgresReady(
	ctx context.Context,
	cfg *config.Postgres,
) (*postgres.Server, error) {
	s, err := postgres.Start(ctx, cfg.DataDir)
	if err != nil {
		return nil, err
	}

	if err := s.EnsureDatabase(
		ctx,
		cfg.Database,
		cfg.Username,
		cfg.Password,
	); err != nil {
		s.Stop(ctx)
		return nil, err
	}

	return s, nil
}

func startMiniflux(
	ctx context.Context,
	baseURL string,
	cfg *config.Info,
) (*miniflux.Server, error) {
	return miniflux.Start(
		ctx,
		miniflux.WithAdmin(
			cfg.Miniflux.AdminUsername,
			cfg.Miniflux.AdminPassword),
		miniflux.WithDatabase(
			cfg.Postgres.Database,
			cfg.Postgres.Username, cfg.Postgres.Password),
		miniflux.WithRunMigrations(true),
		miniflux.WithListenAddress(backendAddr),
		miniflux.WithBaseURL(baseURL),
	)
}

func getDomain(
	ctx context.Context,
	svc *tsweb.Service,
) (string, error) {
	c, err := svc.LocalClient()
	if err != nil {
		return "", err
	}

	var status *ipnstate.Status
	for {
		status, err = c.Status(ctx)
		if err != nil {
			return "", err
		}
		if status.BackendState == "Running" {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	return strings.Trim(status.Self.DNSName, "."), nil
}

func main() {
	var flags Flags
	flags.Register(flag.CommandLine)
	flag.Parse()

	lg := logging.MustSetup()

	var cfg config.Info
	if err := cfg.ReadFile(flags.ConfigFile); err != nil {
		lg.Fatal("unable to read config file",
			zap.Error(err),
			zap.String("config-file", flags.ConfigFile))
		return
	}

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	lg.Info("starting reader",
		zap.String("postgress.data-dir", cfg.Postgres.DataDir))

	svc, err := tsweb.Start(&tsnet.Server{
		AuthKey:  cfg.Tailscale.AuthKey,
		Hostname: cfg.Tailscale.Hostname,
		Dir:      cfg.Tailscale.StateDir,
		Logf: func(format string, args ...any) {
			// lg.Info(fmt.Sprintf(format, args...))
		},
	})
	if err != nil {
		lg.Fatal("unable to start tailscale",
			zap.Error(err))
	}
	defer svc.Close()

	ch := make(chan error, 1)

	go func() {
		ch <- svc.RedirectHTTP(ctx)
	}()

	l, err := svc.ListenTLS("tcp", ":https")
	if err != nil {
		lg.Fatal("unable to listen for https",
			zap.Error(err))
	}
	defer l.Close()

	go func() {
		ch <- web.Serve(ctx, l, "http://"+backendAddr)
	}()

	pg, err := ensurePostgresReady(ctx, &cfg.Postgres)
	if err != nil {
		lg.Fatal("unable to ensure postgres is ready",
			zap.Error(err))
	}
	defer pg.Stop(context.Background())

	lg.Info("postgres started", zap.Int("pid", 0))

	domain, err := getDomain(ctx, svc)
	if err != nil {
		lg.Fatal("unable to get tailscale domain",
			zap.Error(err))
	}
	lg.Info("tailscale domain", zap.String("domain", domain))

	mf, err := startMiniflux(ctx,
		fmt.Sprintf("https://%s/", domain),
		&cfg)
	if err != nil {
		lg.Fatal("unable to start miniflux",
			zap.Error(err))
	}
	defer mf.Stop()

	select {
	case <-ctx.Done():
	case err := <-ch:
		if err != nil {
			lg.Fatal("could not serve http/https",
				zap.Error(err))
		}
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()
}
