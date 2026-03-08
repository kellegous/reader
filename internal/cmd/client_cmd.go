package cmd

import (
	"fmt"
	"net/http"

	"github.com/kellegous/reader/reader_connect"
	"github.com/spf13/cobra"
)

type Codec int

func (c *Codec) Set(v string) error {
	switch v {
	case "protobuf":
		*c = CodecProtobuf
		return nil
	case "json":
		*c = CodecJSON
		return nil
	default:
		return fmt.Errorf("invalid codec: %s", v)
	}
}

func (c Codec) String() string {
	switch c {
	case CodecProtobuf:
		return "protobuf"
	case CodecJSON:
		return "json"
	default:
		return "unknown"
	}
}

func (c Codec) Type() string {
	return "codec"
}

const (
	CodecProtobuf Codec = iota
	CodecJSON
)

type clientFlags struct {
	BaseURL string
}

func (f *clientFlags) NewClient(client *http.Client) reader_connect.ReaderClient {
	return reader_connect.NewReaderClient(client, f.BaseURL)
}

func clientCmd() *cobra.Command {
	var flags clientFlags

	cmd := &cobra.Command{
		Use:   "client",
		Short: "Client for the reader service",
	}

	cmd.Flags().StringVar(
		&flags.BaseURL,
		"base-url",
		"http://localhost:8080/rpc/",
		"Base URL of the reader service",
	)

	cmd.AddCommand(clientCheckHealthCmd(&flags))
	cmd.AddCommand(clientGetEntryTextCmd(&flags))

	return cmd
}
