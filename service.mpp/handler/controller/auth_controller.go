package controller

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/routing"
)

// Environment variables.
var (
	AuthHTTPHost = os.Getenv("AUTH_HTTP_HOST")
)

// AuthController is used to handle HTTP requests to the auth service,
// and to provide a reverse proxy.
type AuthController struct{}

// NewAuthController returns a new instance of AuthController and
// creates a reverse proxy and regsiters a route with the router.
func NewAuthController(r *routing.Router) *AuthController {
	host, _ := url.Parse(AuthHTTPHost)
	proxy := httputil.NewSingleHostReverseProxy(host)

	r.HandleProxy("/api/auth/", http.StripPrefix("/api/auth/", proxy))

	return new(AuthController)
}
