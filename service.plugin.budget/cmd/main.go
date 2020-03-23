package main

import (
	"os"
	"os/signal"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/interface/http"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/interface/rpc"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/registry"
)

func main() {
	bootstrap.Register(os.Getenv("NAME"), os.Getenv("HOSTNAME"))
	defer bootstrap.Unregister(os.Getenv("NAME"))

	ctn := registry.Build()

	web := http.New(ctn)
	go web.Serve()

	server := rpc.Build(ctn)
	go server.Serve()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt)
	<-quit

	go func() {
		<-quit
		server.Shutdown(bootstrap.ShutdownForce)
	}()

	server.Shutdown(bootstrap.ShutdownGraceful)

	ctn.Clean()
}
