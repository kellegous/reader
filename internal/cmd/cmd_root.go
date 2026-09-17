package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reader",
		Short: "Reader is a minflux-based feed reader service",
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
