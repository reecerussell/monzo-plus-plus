package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.auth/interface/http"
	"github.com/reecerussell/monzo-plus-plus/service.auth/interface/rpc"
	"github.com/reecerussell/monzo-plus-plus/service.auth/permission"
	"github.com/reecerussell/monzo-plus-plus/service.auth/registry"
)

func main() {
	ensureDatabaseIsSetup()

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

func ensureDatabaseIsSetup() {
	var err error
	for err != nil {
		select {
		case <-time.After(2 * time.Minute):
			panic(err)
		default:
			db, err := sql.Open("mysql", os.Getenv("CONN_STRING"))
			if err != nil {
				break
			}

			err = db.Ping()
		}
	}
}
