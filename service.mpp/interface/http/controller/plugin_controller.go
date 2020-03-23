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
	PluginsHTTPHost = os.Getenv("PLUGINS_HTTP_HOST")
)

// PluginProxy handles HTTP requests to /api/plugins and forwards the requests to the plugins api.
func PluginProxy() http.HandlerFunc {
	remote, _ := url.Parse(PluginsHTTPHost)
	proxy := httputil.NewSingleHostReverseProxy(remote)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.Replace(r.URL.Path, "/api/plugins", "", 1)

		proxy.ServeHTTP(w, r)
	})
}
