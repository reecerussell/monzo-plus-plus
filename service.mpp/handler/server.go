package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/libraries/routing"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/handler/controller"
)

func Build() *bootstrap.HTTPServer {
	r := routing.NewRouter()

	controller.UI(r)
	_ = controller.NewMonzoController(r)
	_ = controller.NewAuthController(r)
	_ = controller.NewPluginController(r)

	s := bootstrap.BuildServer(&http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")),
	})
	s.CORS(false)

	return s
}
