package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/reecerussell/monzo-plus-plus/service.registry/proto"
	"github.com/reecerussell/monzo-plus-plus/service.registry/service"

	"google.golang.org/grpc"
)

// Environment variable for RPC port
var RPCPort = os.Getenv("RPC_PORT")

func main() {
	server := grpc.NewServer()
	proto.RegisterRegistryServiceServer(server, service.DefaultRegistry)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", RPCPort))
	if err != nil {
		panic(err)
	}

	go func() {
		log.Fatal(server.Serve(lis))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	server.GracefulStop()
}
