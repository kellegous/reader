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
	// NOTE(kellegous):
	// when using a auth_proxy_header, miniflux only uses that to create
	// a new session (if there isn't one already). New sessions are created
	// by the / path. We highjack that path to redirect it to /ui/ which means
	// we break authentication for miniflux. This creates an endpoint that can
	// be called by /ui/ to ensure there is a valid session. The http handler
	// simple requests the / path from miniflux and returns the set-cookie
	// headers that are needed to refresh authentication.
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, beURL.String()+"/", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for k, v := range headers {
			req.Header.Add(k, v)
		}

		// We need a client that doesn't follow redirects because we just need
		// to return the cookies that come back in the first redirect response.
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		res, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		if res.StatusCode >= http.StatusBadRequest {
			http.Error(w, res.Status, res.StatusCode)
			return
		}

		// the response will contain the session and user cookies that
		// are needed to refresh authentication.
		for k, v := range res.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
