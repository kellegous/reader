package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kellegous/glue/devmode"
	"github.com/kellegous/glue/logging"
	"github.com/kellegous/poop"
	"github.com/kellegous/reader"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"miniflux.app/v2/client"

	"github.com/kellegous/reader/internal/config"
	"github.com/kellegous/reader/internal/miniflux"
	"github.com/kellegous/reader/internal/postgres"
	"github.com/kellegous/reader/internal/ui"
	"github.com/kellegous/reader/internal/web"
)

const (
	backendAddr     = "127.0.0.1:9090"
	authProxyHeader = "X-Reader-User"
)

// type DevMode struct {
// 	Root string
// 	Port int
// }

// func (d *DevMode) IsZero() bool {
// 	return d.Port == 0
// }

// func (d *DevMode) Set(v string) error {
// 	root, ps, ok := strings.Cut(v, ":")
// 	if !ok {
// 		root = "."
// 		ps = v
// 	}
// 	port, err := strconv.Atoi(ps)
// 	if err != nil {
// 		return err
// 	}
// 	d.Port = port
// 	d.Root = root
// 	return nil
// }

// func (d *DevMode) String() string {
// 	return fmt.Sprintf("%s:%d", d.Root, d.Port)
// }

// func (d *DevMode) Type() string {
// 	return "root:port"
// }

type serverFlags struct {
	ConfigFile string
	Debug      bool
	DevMode    devmode.Flag
}

func serverCmd() *cobra.Command {
	var flags serverFlags

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the reader server",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runServer(cmd, &flags); err != nil {
				// logging.L(cmd.Context()).Fatal("unable to start server", zap.Error(err))
				poop.HitFan(err)
			}
		},
	}

	cmd.Flags().StringVar(&flags.ConfigFile, "config-file", "reader.yaml", "Path to the config file")
	cmd.Flags().BoolVar(&flags.Debug, "debug", false, "Enable debug logging")
	cmd.Flags().Var(&flags.DevMode, "dev-mode", "Enable dev mode (HMR loading in ui)")
	return cmd
}

func runServer(cmd *cobra.Command, flags *serverFlags) error {
	var cfg config.Info
	if err := cfg.ReadFile(flags.ConfigFile); err != nil {
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
		flags.Debug)
	if err != nil {
		return poop.Chain(err)
	}
	defer mf.Stop()

	ch := make(chan error, 1)

	go func() {
		ch <- poop.Chain(runWeb(ctx, flags, &cfg, mf))
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
	opts := []miniflux.Option{
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
	}

	if cfg.Miniflux.AutoLoginAs != "" {
		opts = append(
			opts, miniflux.WithAuthProxy(
				authProxyHeader,
				true,
				[]string{cfg.Miniflux.AutoLoginAs}))
	}

	s, err := miniflux.Start(ctx, opts...)
	if err != nil {
		return nil, poop.Chain(err)
	}

	if err := s.WaitForReady(ctx, time.Minute); err != nil {
		return nil, poop.Chain(err)
	}

	return s, nil
}

func getAssets(ctx context.Context, devMode *devmode.Flag) (http.Handler, error) {
	if !devMode.IsEnabled() {
		a, err := ui.Assets()
		if err != nil {
			return nil, err
		}
		return http.StripPrefix("/ui/", a), nil
	}

	return devmode.AssetsFromVite(ctx, devMode)
}

func runWeb(
	ctx context.Context,
	flags *serverFlags,
	cfg *config.Info,
	mf *miniflux.Server,
) error {
	l, err := net.Listen("tcp", cfg.Web.Addr)
	if err != nil {
		return poop.Chain(err)
	}
	defer l.Close()

	assets, err := getAssets(ctx, &flags.DevMode)
	if err != nil {
		return poop.Chain(err)
	}

	var mfc *client.Client
	headers := map[string]string{}
	if l := cfg.Miniflux.AutoLoginAs; l != "" {
		headers[authProxyHeader] = l

		key, err := mf.ProvisionUser(ctx, l)
		if err != nil {
			return poop.Chain(err)
		}
		mfc = mf.Client(client.WithAPIKey(key.Token))
	} else {
		mfc = mf.Client()
	}

	if flags.DevMode.IsEnabled() {
		go func() {
			ctx, done := context.WithTimeout(ctx, 10*time.Second)
			defer done()

			if err := flags.DevMode.WaitForReady(ctx); err != nil {
				logging.L(ctx).Fatal("unable to wait for dev mode", zap.Error(err))
			}

			flags.DevMode.PrintBanner(os.Stdout, cfg.Web.Addr)
		}()
	}

	return web.Serve(ctx, l, mf, assets, headers, mfc, &reader.Config{
		Ollama: &reader.Config_Ollama{
			Url:   cfg.Ollama.URL,
			Model: cfg.Ollama.Model,
		},
	})
}
