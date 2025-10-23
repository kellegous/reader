package cmd

import (
	"net/http"
	"time"

	"github.com/kellegous/poop"
	"github.com/kellegous/reader"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func clientGetEntriesByWeekCmd(flags *clientFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "get-entries-by-week",
		Short: "Get entries by week",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runClientGetEntriesByWeek(cmd, flags); err != nil {
				poop.HitFan(err)
			}
		},
	}
}

func runClientGetEntriesByWeek(cmd *cobra.Command, flags *clientFlags) error {
	client, err := flags.NewClient(http.DefaultClient)
	if err != nil {
		return poop.Chain(err)
	}

	// TODO: this needs to be configurable

	now := time.Now()

	res, err := client.GetEntriesByWeek(cmd.Context(), &reader.GetEntriesByWeekRequest{
		StartDate: timestamppb.New(now),
		EndDate:   timestamppb.New(now.AddDate(0, 0, 30)),
	})
	if err != nil {
		return poop.Chain(err)
	}

	m := protojson.MarshalOptions{
		Indent:        "  ",
		UseProtoNames: true,
	}

	b, err := m.Marshal(res)
	if err != nil {
		return poop.Chain(err)
	}

	cmd.Println(string(b))

	return nil
}
