package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/reecerussell/monzo-plus-plus/service.auth/interface/rpc"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.auth/interface/http"
	"github.com/reecerussell/monzo-plus-plus/service.auth/permission"
	"github.com/reecerussell/monzo-plus-plus/service.auth/registry"
)

func main() {
	ctn := registry.Build()

	permission.Build(ctn.Resolve(registry.ServicePermissionsRepository).(repository.PermissionsRepository))

	web := http.Build(ctn)
	go web.Serve()

	gRPC := rpc.Build(ctn)
	go gRPC.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Server shutting down...")

	ctn.Clean()
	web.Shutdown(bootstrap.ShutdownGraceful)
	gRPC.Shutdown(bootstrap.ShutdownGraceful)

	log.Println("Server shutdown!")
}
