package http

import (
	"net/http"

	"github.com/reecerussell/monzo-plus-plus/libraries/sse"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/routing"
)

// Build returns a HTTPServer.
func Build(b *sse.Broker) *bootstrap.HTTPServer {
	r := routing.NewRouter()
	r.Get("/", b)

	return bootstrap.BuildServer(&http.Server{
		Handler: r,
	})
}
