package web

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kellegous/reader"
	"github.com/kellegous/reader/internal/miniflux"
)

func Serve(
	ctx context.Context,
	l net.Listener,
	ms *miniflux.Server,
) error {
	beURL, err := url.Parse(ms.BaseURL())
	if err != nil {
		return err
	}

	m := http.NewServeMux()

	m.Handle("/", httputil.NewSingleHostReverseProxy(beURL))
	m.Handle(reader.ReaderPathPrefix, reader.NewReaderServer(&rpc{
		client: ms.Client(),
	}))

	return http.Serve(l, m)
}
