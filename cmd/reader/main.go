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
	}

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	lg.Info("starting reader",
		zap.String("postgress.data-dir", cfg.Postgres.DataDir))

	// TODO(knorton): Join tailnet

	pg, err := postgres.Start(ctx, cfg.Postgres.DataDir)
	if err != nil {
		lg.Fatal("unable to start postgres",
			zap.Error(err))
	}

	lg.Info("postgres started",
		zap.Int("pid", pg.Process().Pid))

	// TODO(knorton): Start miniflux

	<-ctx.Done()
}
