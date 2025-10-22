package cmd

import (
	"github.com/kellegous/glue/logging"
	"github.com/kellegous/poop"
	"github.com/kellegous/reader/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

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

	return nil
}
