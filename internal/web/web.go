package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kellegous/reader"
	"github.com/kellegous/reader/internal/config"
	"github.com/kellegous/reader/internal/miniflux"
	"miniflux.app/v2/client"
)

func Serve(
	ctx context.Context,
	cfg *config.Info,
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
	var api *client.Client
	if l := cfg.Miniflux.AutoLoginAs; l != "" {
		key := ms.APIKeyFor(l)
		if key == nil {
			return fmt.Errorf("no api key for user: %s", l)
		}
		api = ms.Client(client.WithAPIKey(key.Token))
	} else {
		api = ms.Client()
	}

	m.Handle("/", p)
	m.Handle(reader.ReaderPathPrefix, reader.NewReaderServer(&rpc{
		client: api,
	}))
	m.Handle("/ui/", assets)

	return http.Serve(l, m)
}
