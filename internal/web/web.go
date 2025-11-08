package web

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kellegous/poop"
	"miniflux.app/v2/client"

	"github.com/kellegous/reader"
	"github.com/kellegous/reader/internal/miniflux"
)

func getMinifluxClient(
	ctx context.Context,
	ms *miniflux.Server,
	username string,
) (*client.Client, error) {
	if username == "" {
		return ms.Client(), nil
	}

	key, err := ms.ProvisionUser(ctx, username)
	if err != nil {
		return nil, poop.Chain(err)
	}

	return ms.Client(client.WithAPIKey(key.Token)), nil
}

// TODO(kellegous): consolidate these args into a single options
// struct.
func Serve(
	ctx context.Context,
	l net.Listener,
	ms *miniflux.Server,
	assets http.Handler,
	headers map[string]string,
	username string,
	cfg *reader.Config,
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

	api, err := getMinifluxClient(ctx, ms, username)
	if err != nil {
		return poop.Chain(err)
	}

	m.Handle("/", p)
	m.Handle(reader.ReaderPathPrefix, reader.NewReaderServer(&rpc{
		client: api,
		cfg:    cfg,
	}))
	m.Handle("/ui/", assets)

	return http.Serve(l, m)
}
