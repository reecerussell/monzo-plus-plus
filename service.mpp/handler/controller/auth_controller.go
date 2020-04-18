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

type AuthController struct{}

func NewAuthController(r *routing.Router) *AuthController {
	host, _ := url.Parse(AuthHTTPHost)
	proxy := httputil.NewSingleHostReverseProxy(host)

	r.HandleProxy("/auth/", http.StripPrefix("/auth/", proxy))

	return new(AuthController)
}
