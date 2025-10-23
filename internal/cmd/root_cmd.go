package cmd

import (
	"os"

	"github.com/kellegous/glue/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func rootCmd() *cobra.Command {
	var lg *zap.Logger
	cmd := &cobra.Command{
		Use:   "reader",
		Short: "Reader is a minflux-based feed reader service",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			lg = logging.MustSetup()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			lg.Sync()
		},
	}

	cmd.AddCommand(serverCmd())
	cmd.AddCommand(clientCmd())

	return cmd
}

func Execute() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
