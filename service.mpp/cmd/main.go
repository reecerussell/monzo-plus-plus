package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"

	"github.com/reecerussell/monzo-plus-plus/service.mpp/interface/http"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/plugin"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/registry"

	// plugins
	_ "github.com/reecerussell/monzo-plus-plus/service.mpp/plugin/budget"
)

func main() {
	ctn := registry.Build()
	plugin.Build(ctn)

	s := http.Build(ctn)
	go s.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Server shutting down...")

	ctn.Clean()
	s.Shutdown(bootstrap.ShutdownGraceful)

	log.Println("Server shutdown!")
}
