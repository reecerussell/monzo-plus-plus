package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.auth/interface/http/controller"
	"github.com/reecerussell/monzo-plus-plus/service.auth/interface/http/middleware"
)

func Build(ctn *di.Container) *bootstrap.HTTPServer {
	am := middleware.NewAuthenticationMiddleware(ctn)
	r := mux.NewRouter()

	_ = controller.NewTokenController(ctn, r)
	_ = controller.NewUserController(ctn, r)
	_ = controller.NewRoleController(ctn, r)
	_ = controller.NewPermissionsController(ctn, r)

	return bootstrap.BuildServer(&http.Server{
		Handler: am.Handler(r),
		Addr:    fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")),
	})
}
