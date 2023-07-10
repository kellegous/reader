package web

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Serve(
	ctx context.Context,
	l net.Listener,
	be string,
) error {
	beURL, err := url.Parse(be)
	if err != nil {
		return err
	}

	m := http.NewServeMux()

	m.Handle("/", httputil.NewSingleHostReverseProxy(beURL))

	return http.Serve(l, m)
}
