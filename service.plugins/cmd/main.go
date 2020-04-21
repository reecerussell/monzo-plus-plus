package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/interface/http"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/registry"
)

func main() {
	ctn := registry.Build()

	web := http.Build(ctn)
	go web.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Server shutting down...")

	ctn.Clean()
	web.Shutdown(bootstrap.ShutdownGraceful)

	log.Println("Server shutdown!")
}
