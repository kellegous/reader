package cmd

import (
	"fmt"
	"net/http"

	"github.com/kellegous/reader"
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
	Codec   Codec
}

func (f *clientFlags) NewClient(client *http.Client) (reader.Reader, error) {
	switch f.Codec {
	case CodecProtobuf:
		return reader.NewReaderProtobufClient(f.BaseURL, client), nil
	case CodecJSON:
		return reader.NewReaderJSONClient(f.BaseURL, client), nil
	}
	return nil, fmt.Errorf("invalid codec: %s", f.Codec)
}

func clientCmd() *cobra.Command {
	flags := clientFlags{
		Codec: CodecProtobuf,
	}

	cmd := &cobra.Command{
		Use:   "client",
		Short: "Client for the reader service",
	}

	cmd.Flags().StringVar(&flags.BaseURL, "base-url", "http://localhost:8080", "Base URL of the reader service")
	cmd.Flags().Var(&flags.Codec, "codec", "Codec to use for the client")

	cmd.AddCommand(clientCheckHealthCmd(&flags))
	cmd.AddCommand(clientGetEntriesByWeekCmd(&flags))

	return cmd
}
