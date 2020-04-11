package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/routing"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/interface/http/controller"
)

func Build(ctn *di.Container) *bootstrap.HTTPServer {
	r := routing.NewRouter()

	_ = controller.NewMonzoController(ctn, r)
	_ = controller.NewAuthController(r)
	_ = controller.NewPluginController(r)

	return bootstrap.BuildServer(&http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")),
	})
}
