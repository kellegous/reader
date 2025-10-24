package cmd

import (
	"fmt"
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

	res, err := client.GetEntriesByWeek(cmd.Context(), &reader.GetEntriesByWeekRequest{
		Range: &reader.GetEntriesByWeekRequest_NWeeksFromWeek_{
			NWeeksFromWeek: &reader.GetEntriesByWeekRequest_NWeeksFromWeek{
				FromWeekOf: timestamppb.New(time.Now()),
				NWeeks:     5,
			},
		},
		WeekStartsDay: reader.GetEntriesByWeekRequest_MONDAY,
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

	fmt.Println(string(b))

	return nil
}
