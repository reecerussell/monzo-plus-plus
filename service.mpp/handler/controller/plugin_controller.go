package controller

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/routing"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
)

// Environment variables.
var (
	PluginsHTTPHost = os.Getenv("PLUGINS_HTTP_HOST")
)

// PluginController is used to handle HTTP requests and manage
// reverse proxies for the plugin APIs.
type PluginController struct {
	hosts   map[string]*httputil.ReverseProxy
	plugins *httputil.ReverseProxy
}

// NewPluginController returns a new instance of PluginController and
// instantiates the plugin reverse proxy and routes.
func NewPluginController(r *routing.Router) *PluginController {
	pluginsURL, _ := url.Parse(PluginsHTTPHost)

	c := &PluginController{
		hosts:   make(map[string]*httputil.ReverseProxy),
		plugins: httputil.NewSingleHostReverseProxy(pluginsURL),
	}

	r.HandleProxy("/api/plugins/", http.StripPrefix("/api/plugins/", c.plugins))
	r.HandleProxy("/api/plugin/{name}/", http.Handler(http.HandlerFunc(c.HandlePlugin)))

	return c
}

// HandlePlugin acts as a reverse proxy to internal plugin APIs.
func (c *PluginController) HandlePlugin(w http.ResponseWriter, r *http.Request) {
	name := routing.Vars(r)["name"]

	proxy, ok := c.hosts[name]
	if !ok {
		host, err := bootstrap.GetHost(name)
		if err != nil {
			errors.HandleHTTPError(w, r, errors.InternalError(err))
			return
		}

		if host == "" {
			errors.HandleHTTPError(w, r, errors.NotFound("host is empty"))
			return
		}

		url, _ := url.Parse(fmt.Sprintf("http://%s", host))
		c.hosts[name] = httputil.NewSingleHostReverseProxy(url)
		proxy = c.hosts[name]
	}

	handler := http.StripPrefix(fmt.Sprintf("/api/plugin/%s", name), proxy)
	handler.ServeHTTP(w, r)
}
