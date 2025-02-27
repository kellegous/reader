package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	"github.com/kellegous/reader/internal/config"
	"github.com/kellegous/reader/internal/logging"
	"github.com/kellegous/reader/internal/miniflux"
	"github.com/kellegous/reader/internal/postgres"
	"github.com/kellegous/reader/internal/web"
)

const (
	backendAddr = "127.0.0.1:9090"
)

type Flags struct {
	ConfigFile string
	Debug      bool
}

func (f *Flags) Register(fs *flag.FlagSet) {
	fs.StringVar(
		&f.ConfigFile,
		"config-file",
		"reader.yaml",
		"Path to the config file",
	)
	fs.BoolVar(
		&f.Debug,
		"debug",
		false,
		"Enable debug logging",
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
	debug bool,
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
		miniflux.WithDebugLogging(debug),
	)
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

	ch := make(chan error, 1)

	l, err := net.Listen("tcp", cfg.Web.Addr)
	if err != nil {
		lg.Fatal("unable to listen for http",
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

	mf, err := startMiniflux(ctx,
		fmt.Sprintf("https://%s/", cfg.Web.Hostname),
		&cfg,
		flags.Debug)
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
