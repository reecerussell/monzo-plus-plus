package routing

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/reecerussell/monzo-plus-plus/libraries/util"
)

var (
	routeContextKey = util.ContextKey("route")
)

// Route is used to map a URL path to a HTTP handler.
type Route struct {
	Template string
	Handler  http.Handler
}

// getVars scans a given path and matches variables in the template.
func (r *Route) getVars(path string) map[string]string {
	tp := strings.Split(r.Template, "/")[1:]
	up := strings.Split(path, "/")[1:]

	vars := make(map[string]string)
	re := regexp.MustCompile("\\{(.*?)\\}")

	for i, tp := range tp {
		if !re.MatchString(tp) {
			continue
		}

		n := tp[1 : len(tp)-1]
		vars[n] = up[i]
	}

	return vars
}

// Router is used to handle incoming HTTP requests and map them to
// a matching route.
type Router struct {
	mu     *sync.RWMutex
	routes map[string][]*Route
}

// NewRouter returns a new instance of Router.
func NewRouter() *Router {
	return &Router{
		mu:     &sync.RWMutex{},
		routes: make(map[string][]*Route),
	}
}

// Get creates a new route for handling GET requests.
func (rtr *Router) Get(tpl string, hlr http.Handler) *Route {
	return rtr.newRoute(tpl, hlr, http.MethodGet)
}

// GetFunc creates a new route for handling GET requests.
func (rtr *Router) GetFunc(tpl string, hlr http.HandlerFunc) *Route {
	return rtr.newRoute(tpl, http.Handler(hlr), http.MethodGet)
}

// Post creates a new route for handling POST requests.
func (rtr *Router) Post(tpl string, hlr http.Handler) *Route {
	return rtr.newRoute(tpl, hlr, http.MethodPost)
}

// PostFunc creates a new route for handling POST requests.
func (rtr *Router) PostFunc(tpl string, hlr http.HandlerFunc) *Route {
	return rtr.newRoute(tpl, http.Handler(hlr), http.MethodPost)
}

// Put creates a new route for handling PUT requests.
func (rtr *Router) Put(tpl string, hlr http.Handler) *Route {
	return rtr.newRoute(tpl, hlr, http.MethodPut)
}

// PutFunc creates a new route for handling PUT requests.
func (rtr *Router) PutFunc(tpl string, hlr http.HandlerFunc) *Route {
	return rtr.newRoute(tpl, http.Handler(hlr), http.MethodPut)
}

// Delete creates a new route for handling DELETE requests.
func (rtr *Router) Delete(tpl string, hlr http.Handler) *Route {
	return rtr.newRoute(tpl, hlr, http.MethodDelete)
}

// DeleteFunc creates a new route for handling DELETE requests.
func (rtr *Router) DeleteFunc(tpl string, hlr http.HandlerFunc) *Route {
	return rtr.newRoute(tpl, http.Handler(hlr), http.MethodDelete)
}

// Handle creates a new route for handling HTTP requests with the given methods.
func (rtr *Router) Handle(tpl string, hlr http.Handler, mtds ...string) *Route {
	return rtr.newRoute(tpl, hlr, mtds...)
}

// HandleFunc creates a new route for handling HTTP requests with the given methods.
func (rtr *Router) HandleFunc(tpl string, hlr http.HandlerFunc, mtds ...string) *Route {
	return rtr.newRoute(tpl, http.Handler(hlr), mtds...)
}

// newRoute adds a new route to the router.
func (rtr *Router) newRoute(tpl string, hlr http.Handler, mtds ...string) *Route {
	rtr.mu.RLock()
	defer rtr.mu.RUnlock()

	r := &Route{
		Template: tpl,
		Handler:  hlr,
	}

	for _, method := range mtds {
		rts, ok := rtr.routes[method]
		if ok {
			rts = append(rts, r)
			rtr.routes[method] = rts
		} else {
			rts = []*Route{r}
			rtr.routes[method] = rts
		}
	}

	return r
}

// match is used to determine which route handler should be used for a given path and method.
func (rtr *Router) match(path, method string) *Route {
	for _, route := range rtr.routes[method] {
		if strings.HasPrefix(path, route.Template) {
			return route
		}

		tpl := route.Template
		re := regexp.MustCompile("\\{(.*?)\\}")
		tpl = re.ReplaceAllString(tpl, "(.*?)")
		re = regexp.MustCompile(tpl)

		if re.MatchString(path) {
			return route
		}
	}

	return nil
}

// ServeHTTP is a http.Handler function which is used to map incoming
// requests to the appropriate route and handler.
func (rtr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr.mu.Lock()
	defer rtr.mu.Unlock()

	route := rtr.match(r.URL.Path, r.Method)
	if route == nil {
		log.Printf("[ROUTER]: no '%s' route found for '%s'", r.Method, r.URL.Path)
		http.NotFound(w, r)
		return
	}

	ctx := context.WithValue(r.Context(), routeContextKey, route)
	r = r.WithContext(ctx)

	route.Handler.ServeHTTP(w, r)
}
