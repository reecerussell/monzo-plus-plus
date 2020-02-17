package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Environemt variables.
var (
	HTTPPort = os.Getenv("HTTP_PORT")
)

// Shutdown modes.
const (
	ShutdownGraceful = 1 << iota
	ShutdownForce
)

type HTTPServer struct {
	base     *http.Server
	shutdown chan int
}

func BuildServer(s *http.Server) *HTTPServer {
	return &HTTPServer{
		base:     s,
		shutdown: make(chan int),
	}
}

func (hs *HTTPServer) Serve() {
	hs.base.Addr = fmt.Sprintf(":%s", HTTPPort)

	sc := make(chan struct{})
	go hs.listenForShutdown()

	log.Printf("HTTP Server listening on: %s\n", hs.base.Addr)

	if err := hs.base.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("http: serve: %v", err)
	}

	<-sc
}

// Shutdown closes the HTTP server, depending on the mode.
func (hs *HTTPServer) Shutdown(mode int) {
	hs.shutdown <- mode
}

func (hs *HTTPServer) listenForShutdown() {
	go func() {
		mode := <-hs.shutdown
		var err error

		switch mode {
		case ShutdownGraceful:
			log.Printf("HTTP Server gracefully shutting down...")
			err = hs.base.Shutdown(context.Background())
			break
		case ShutdownForce:
			log.Printf("HTTP Server forcefully shutting down...")
			err = hs.base.Close()
			break
		default:
			log.Panicf("http: shutdown: %d is not a valid mode", mode)
			break
		}

		if err != nil {
			log.Fatalf("HTTP Server failed to shutdown: %v", err)
		}

		log.Printf("HTTP Server shutdown.")

		close(hs.shutdown)
	}()
}
