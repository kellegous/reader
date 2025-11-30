package web

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"miniflux.app/v2/client"

	"github.com/kellegous/reader"
	"github.com/kellegous/reader/internal/miniflux"
)

// TODO(kellegous): consolidate these args into a single options
// struct.
func Serve(
	ctx context.Context,
	l net.Listener,
	ms *miniflux.Server,
	assets http.Handler,
	headers map[string]string,
	api *client.Client,
	cfg *reader.Config,
) error {
	beURL, err := url.Parse(ms.BaseURL())
	if err != nil {
		return err
	}

	m := http.NewServeMux()

	m.Handle("/", newMinifluxProxy(beURL, headers))
	m.Handle(reader.ReaderPathPrefix, reader.NewReaderServer(&rpc{
		client: api,
		cfg:    cfg,
	}))
	m.Handle("/refresh-session", newSessionRefresher(beURL, headers))
	m.Handle("/ui/", assets)

	return http.Serve(l, m)
}

func newMinifluxProxy(beURL *url.URL, headers map[string]string) http.Handler {
	p := httputil.NewSingleHostReverseProxy(beURL)
	dir := p.Director
	p.Director = func(r *http.Request) {
		dir(r)
		for k, v := range headers {
			r.Header.Add(k, v)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/ui/", http.StatusTemporaryRedirect)
			return
		}
		p.ServeHTTP(w, r)
	})
}

func newSessionRefresher(beURL *url.URL, headers map[string]string) http.Handler {
	director := func(r *http.Request) {
		for k, v := range headers {
			r.Header.Add(k, v)
		}
		r.URL = beURL
	}
	return &httputil.ReverseProxy{Director: director}
}
