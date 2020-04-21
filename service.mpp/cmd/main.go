package main

import (
	"os"
	"os/signal"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/handler"
)

func main() {
	s := handler.Build()
	go s.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	s.Shutdown(bootstrap.ShutdownGraceful)
}
