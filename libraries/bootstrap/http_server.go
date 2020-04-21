package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
)

// Environemt variables.
var (
	HTTPPort    = os.Getenv("HTTP_PORT")
	HTTPSPort   = os.Getenv("HTTPS_PORT")
	SSLCertFile = os.Getenv("SSL_CERT_FILE")
	SSLKeyFile  = os.Getenv("SSL_KEY_FILE")
)

// Shutdown modes.
const (
	ShutdownGraceful = 1 << iota
	ShutdownForce
)

type HTTPServer struct {
	base     *http.Server
	cors     bool
	shutdown chan int
}

func BuildServer(s *http.Server) *HTTPServer {
	return &HTTPServer{
		base:     s,
		cors:     true,
		shutdown: make(chan int),
	}
}

func (hs *HTTPServer) CORS(enabled bool) {
	hs.cors = enabled
}

func (hs *HTTPServer) Serve() {
	hs.base.Addr = fmt.Sprintf(":%s", HTTPPort)

	if hs.cors {
		hs.base.Handler = panicHandler(corsHandler(hs.base.Handler))
	} else {
		hs.base.Handler = panicHandler(hs.base.Handler)
	}

	sc := make(chan struct{})
	go hs.listenForShutdown()

	log.Printf("HTTP Server listening on: %s\n", hs.base.Addr)

	if err := hs.base.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("http: serve: %v", err)
	}

	<-sc
}

func (hs *HTTPServer) ServeTLS() {
	hs.base.Addr = fmt.Sprintf(":%s", HTTPSPort)

	if hs.cors {
		hs.base.Handler = panicHandler(corsHandler(hs.base.Handler))
	} else {
		hs.base.Handler = panicHandler(hs.base.Handler)
	}

	sc := make(chan struct{})
	go hs.listenForShutdown()

	log.Printf("HTTP Server listening on: %s\n", hs.base.Addr)

	if err := hs.base.ListenAndServeTLS(SSLCertFile, SSLKeyFile); err != http.ErrServerClosed {
		log.Fatalf("http: serve: %v", err)
	}

	<-sc
}

// Shutdown closes the HTTP server, depending on the mode.
func (hs *HTTPServer) Shutdown(mode int) {
	hs.shutdown <- mode
}

func (hs *HTTPServer) listenForShutdown() {
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
}

func panicHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			p := recover()
			if p == nil {
				return
			}

			err := fmt.Errorf("%v", p)
			errors.HandleHTTPError(w, r, errors.InternalError(err))
		}()

		h.ServeHTTP(w, r)
	})
}

func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)

		if _, ok := w.Header()["Access-Control-Allow-Origin"]; !ok {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}

		if _, ok := w.Header()["Access-Control-Allow-Methods"]; !ok {
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
		}

		if _, ok := w.Header()["Access-Control-Allow-Headers"]; !ok {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
	})
}
