package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.auth/interface/http/controller"
)

func Build(ctn *di.Container) *bootstrap.HTTPServer {
	r := mux.NewRouter()

	_ = controller.NewTokenController(ctn, r)
	_ = controller.NewUserController(ctn, r)
	_ = controller.NewUserRoleController(ctn, r)
	_ = controller.NewRoleController(ctn, r)

	return bootstrap.BuildServer(&http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")),
	})
}
