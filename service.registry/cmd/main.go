package main

import (
	"os"
	"os/signal"

	"github.com/reecerussell/monzo-plus-plus/libraries/bootstrap"

	"github.com/reecerussell/monzo-plus-plus/service.registry/proto"
	"github.com/reecerussell/monzo-plus-plus/service.registry/service"

	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	proto.RegisterRegistryServiceServer(server, service.DefaultRegistry)

	rpc := bootstrap.BuildRPCServer(server)
	go rpc.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	server.GracefulStop()
}
