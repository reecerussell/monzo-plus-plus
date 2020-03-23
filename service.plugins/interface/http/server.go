package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/permission"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/interface/http/controller"
)

// Build returns a HTTPServer with plugin routes and authentication middleware.
func Build(ctn *di.Container) *bootstrap.HTTPServer {
	r := mux.NewRouter()

	_ = controller.NewPluginController(ctn, r)

	return bootstrap.BuildServer(&http.Server{
		Handler: permission.Middleware(r, permission.PermissionPluginManager),
		Addr:    fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")),
	})
}
