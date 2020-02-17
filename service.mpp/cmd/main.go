package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/reecerussell/monzo-plus-plus/service.mpp/interface/http"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/plugin"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/registry"

	// plugins
	_ "github.com/reecerussell/monzo-plus-plus/service.mpp/plugin/budget"
)

func main() {
	ctn := registry.Build()
	plugin.Build(ctn)

	s := http.NewServer(ctn)
	go s.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Server shutting down...")

	ctn.Clean()
	s.Shutdown()

	log.Println("Server shutdown!")
}
