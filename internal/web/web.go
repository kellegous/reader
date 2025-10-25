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
	assets http.Handler,
	headers map[string]string,
) error {
	beURL, err := url.Parse(ms.BaseURL())
	if err != nil {
		return err
	}

	m := http.NewServeMux()

	p := httputil.NewSingleHostReverseProxy(beURL)
	dir := p.Director
	p.Director = func(r *http.Request) {
		dir(r)
		for k, v := range headers {
			r.Header.Add(k, v)
		}
	}

	m.Handle("/", p)
	m.Handle(reader.ReaderPathPrefix, reader.NewReaderServer(&rpc{
		client: ms.Client(),
	}))
	m.Handle("/ui/", assets)

	return http.Serve(l, m)
}
