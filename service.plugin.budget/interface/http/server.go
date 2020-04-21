package http

import (
	"net/http"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/permission"
	"github.com/reecerussell/monzo-plus-plus/libraries/routing"

	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/interface/http/controller"
)

func New(ctn *di.Container) *bootstrap.HTTPServer {
	r := routing.NewRouter()

	_ = controller.NewPreferencesController(ctn, r)

	return bootstrap.BuildServer(&http.Server{
		Handler: permission.Middleware(r),
	})
}
