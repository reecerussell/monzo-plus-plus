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
	err := bootstrap.Register(os.Getenv("NAME"), os.Getenv("HOSTNAME"))
	if err != nil {
		panic(err)
	}
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
		web.Shutdown(bootstrap.ShutdownForce)
		server.Shutdown(bootstrap.ShutdownForce)
	}()

	web.Shutdown(bootstrap.ShutdownGraceful)
	server.Shutdown(bootstrap.ShutdownGraceful)

	ctn.Clean()
}
