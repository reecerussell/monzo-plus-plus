package controller

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
)

// Environment variables.
var (
	PluginsHTTPHost = os.Getenv("PLUGINS_HTTP_HOST")
)

type PluginController struct {
	hosts   map[string]*httputil.ReverseProxy
	plugins *httputil.ReverseProxy
}

func NewPluginController(r *mux.Router) *PluginController {
	pluginsURL, _ := url.Parse(PluginsHTTPHost)

	c := &PluginController{
		hosts:   make(map[string]*httputil.ReverseProxy),
		plugins: httputil.NewSingleHostReverseProxy(pluginsURL),
	}

	r.HandleFunc("/api/plugins/", c.HandlePlugins)
	r.HandleFunc("/api/plguin/{name}", c.HandlePluginAPI)

	return c
}

// HandlePlugins handles HTTP requests to /api/plugins and forwards the requests to the plugins api.
func (c *PluginController) HandlePlugins(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.Replace(r.URL.Path, "/api/plugins", "", 1)

	c.plugins.ServeHTTP(w, r)
}

// HandlePluginAPI acts as a reverse proxy to internal plugin APIs.
func (c *PluginController) HandlePluginAPI(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		http.NotFound(w, r)
		return
	}

	proxy, ok := c.hosts[name]
	if !ok {
		host, err := bootstrap.GetHost(name)
		if err != nil {
			errors.HandleHTTPError(w, r, errors.InternalError(err))
			return
		}

		if host == "" {
			http.NotFound(w, r)
			return
		}

		url, _ := url.Parse(host)
		c.hosts[name] = httputil.NewSingleHostReverseProxy(url)
		proxy = c.hosts[name]
	}

	r.URL.Path = strings.Replace(r.URL.Path, fmt.Sprintf("/api/plugin/%s", name), "", 1)
	proxy.ServeHTTP(w, r)
}
