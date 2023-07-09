package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/kellegous/reader/pkg/config"
	"github.com/kellegous/reader/pkg/logging"
	"github.com/kellegous/reader/pkg/postgres"
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

	// TODO(knorton): Join tailnet

	pg, err := ensurePostgresReady(ctx, &cfg.Postgres)
	if err != nil {
		lg.Fatal("unable to ensure postgres is ready",
			zap.Error(err))
		return
	}
	defer pg.Stop(context.Background())

	lg.Info("postgres started", zap.Int("pid", 0))

	// TODO(knorton): Start miniflux

	<-ctx.Done()
}
