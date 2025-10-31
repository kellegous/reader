package cmd

import (
	"net/http"
	"strconv"

	"github.com/kellegous/poop"
	"github.com/kellegous/reader"
	"github.com/spf13/cobra"
)

func clientGetEntryTextCmd(flags *clientFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "get-entry-text",
		Short: "Get the text of an entry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := runClientGetEntryText(cmd, flags, args[0]); err != nil {
				poop.HitFan(err)
			}
		},
	}
}

func runClientGetEntryText(cmd *cobra.Command, flags *clientFlags, idStr string) error {
	client, err := flags.NewClient(http.DefaultClient)
	if err != nil {
		return poop.Chain(err)
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return poop.Chain(err)
	}

	entry, err := client.GetEntryText(cmd.Context(), &reader.GetEntryTextRequest{
		EntryId: id,
	})
	if err != nil {
		return poop.Chain(err)
	}

	cmd.Println(entry.Text)

	return nil
}
