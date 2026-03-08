package cmd

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

func clientCheckHealthCmd(flags *clientFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "check-health",
		Short: "Check the health of the reader service",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runClientCheckHealth(cmd, flags); err != nil {
				poop.HitFan(err)
			}
		},
	}
}

func runClientCheckHealth(cmd *cobra.Command, flags *clientFlags) error {
	client := flags.NewClient(http.DefaultClient)

	if _, err := client.CheckHealth(cmd.Context(), connect.NewRequest(&emptypb.Empty{})); err != nil {
		return poop.Chain(err)
	}

	cmd.Println("👍")

	return nil
}
