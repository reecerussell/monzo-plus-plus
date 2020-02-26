package controller

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// Environment variables.
var (
	AuthHTTPHost = os.Getenv("AUTH_HTTP_HOST")
)

// AuthProxy is a http.HandlerFunc that returns a handler to reverse proxy to the auth service.
func AuthProxy() http.HandlerFunc {
	remote, _ := url.Parse(AuthHTTPHost)
	proxy := httputil.NewSingleHostReverseProxy(remote)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Replace(r.URL.Path, "/auth", "", 1)

		proxy.ServeHTTP(w, r)
	})
}
