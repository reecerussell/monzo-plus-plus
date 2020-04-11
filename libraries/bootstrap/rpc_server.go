package bootstrap

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

// Environment variables.
var (
	RPCPort = os.Getenv("RPC_PORT")
)

type RPCServer struct {
	base     *grpc.Server
	shutdown chan int
}

func BuildRPCServer(b *grpc.Server) *RPCServer {
	return &RPCServer{
		base:     b,
		shutdown: make(chan int),
	}
}

func (rs *RPCServer) Serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", RPCPort))
	if err != nil {
		log.Panicf("rpc: failed to listen: %v", err)
	}

	s := make(chan struct{})
	go rs.listenForShutdown()

	log.Printf("RPC Server listening on: %s\n", lis.Addr().String())
	if err := rs.base.Serve(lis); err != nil {
		log.Fatalf("rpc: serve: %v", err)
	}

	<-s
}

// Shutdown closes the RPC server, depending on the mode.
func (rs *RPCServer) Shutdown(mode int) {
	rs.shutdown <- mode
}

func (rs *RPCServer) listenForShutdown() {
	go func() {
		mode := <-rs.shutdown

		switch mode {
		case ShutdownGraceful:
			log.Printf("RPC Server gracefully shutting down...")
			rs.base.GracefulStop()
			break
		case ShutdownForce:
			log.Printf("RPC Server forcefully shutting down...")
			rs.base.Stop()
			break
		default:
			log.Panicf("rpc: shutdown: %d is not a valid mode", mode)
			break
		}

		log.Printf("RPC Server shutdown.")

		close(rs.shutdown)
	}()
}
