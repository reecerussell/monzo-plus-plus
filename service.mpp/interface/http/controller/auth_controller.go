package controller

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// Environment variables.
var (
	AuthHTTPHost = os.Getenv("AUTH_HTTP_HOST")
)

type AuthController struct {
	auth *httputil.ReverseProxy
}

func NewAuthController(r *mux.Router) *AuthController {
	host, _ := url.Parse(AuthHTTPHost)

	c := &AuthController{
		auth: httputil.NewSingleHostReverseProxy(host),
	}

	r.HandleFunc("/auth/", c.HandleAuthAPI)

	return c
}

// HandleAuthAPI handles request to the Auth service using a reverse proxy.
func (c *AuthController) HandleAuthAPI(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.Replace(r.URL.Path, "/auth", "", 1)

	c.auth.ServeHTTP(w, r)
}
