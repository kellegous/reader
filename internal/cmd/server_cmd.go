package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/kellegous/glue/logging"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/kellegous/reader/internal/config"
	"github.com/kellegous/reader/internal/miniflux"
	"github.com/kellegous/reader/internal/postgres"
	"github.com/kellegous/reader/internal/web"
)

const backendAddr = "127.0.0.1:9090"

func serverCmd() *cobra.Command {
	var flags struct {
		ConfigFile string
		Debug      bool
	}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the reader server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runServer(cmd, flags.ConfigFile, flags.Debug); err != nil {
				logging.L(cmd.Context()).Fatal("unable to start server", zap.Error(err))
			}
		},
	}

	cmd.Flags().StringVar(&flags.ConfigFile, "config-file", "reader.yaml", "Path to the config file")
	cmd.Flags().BoolVar(&flags.Debug, "debug", false, "Enable debug logging")
	return cmd
}

func runServer(cmd *cobra.Command, configFile string, debug bool) error {
	var cfg config.Info
	if err := cfg.ReadFile(configFile); err != nil {
		return poop.Chain(err)
	}

	ctx, done := signal.NotifyContext(cmd.Context(), os.Interrupt)
	defer done()

	lg := logging.L(cmd.Context())

	lg.Info("starting reader",
		zap.String("postgress.data-dir", cfg.Postgres.DataDir))

	pg, err := ensurePostgresReady(ctx, &cfg.Postgres)
	if err != nil {
		return poop.Chain(err)
	}
	defer pg.Stop(context.Background())

	// TODO(knorton): get pid from postgres
	lg.Info("postgres started", zap.Int("pid", 0))

	mf, err := startMiniflux(
		ctx,
		fmt.Sprintf("https://%s/", cfg.Web.Hostname),
		&cfg,
		debug)
	if err != nil {
		return poop.Chain(err)
	}
	defer mf.Stop()

	ch := make(chan error, 1)

	l, err := net.Listen("tcp", cfg.Web.Addr)
	if err != nil {
		return poop.Chain(err)
	}
	defer l.Close()

	go func() {
		ch <- web.Serve(ctx, l, mf)
	}()

	select {
	case <-ctx.Done():
	case err := <-ch:
		if err != nil {
			return poop.Chain(err)
		}
	}

	return nil
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
