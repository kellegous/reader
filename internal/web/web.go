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

	p := httputil.NewSingleHostReverseProxy(beURL)
	dir := p.Director
	p.Director = func(r *http.Request) {
		for k, v := range headers {
			r.Header.Add(k, v)
		}
		dir(r)
	}

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// TODO(knorton): This doesn't work because miniflux redirects
		// from /unread to / and back to /unread.
		// if r.URL.Path == "/" {
		// 	http.Redirect(w, r, "/ui/", http.StatusTemporaryRedirect)
		// 	return
		// }
		p.ServeHTTP(w, r)
	})
	m.Handle(reader.ReaderPathPrefix, reader.NewReaderServer(&rpc{
		client: api,
		cfg:    cfg,
	}))
	m.Handle("/ui/", assets)

	return http.Serve(l, m)
}
