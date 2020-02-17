package http

import (
	"net/http"

	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/interface/http/controller"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
)

func New(ctn *di.Container) *bootstrap.HTTPServer {
	mux := http.NewServeMux()

	_ = controller.NewPreferencesController(ctn, mux)

	return bootstrap.BuildServer(&http.Server{
		Handler: mux,
	})
}
