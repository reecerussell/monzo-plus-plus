package http

import (
	"context"
	"fmt"
	"log"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/interface/http/controller"

	"net/http"
	"os"
)

// Environemt variables
var (
	HTTPPort = os.Getenv("HTTP_PORT")
)

type Server struct {
	mux *http.ServeMux
	s   *http.Server
}

func NewServer(ctn *di.Container) *Server {
	mux := &http.ServeMux{}

	controller.NewMonzoController().Apply(ctn, mux)

	return &Server{
		mux: mux,
		s:   &http.Server{},
	}
}

func (s *Server) Serve() {
	// mux := &http.ServeMux{}
	// mux.Handle("/", panicHandler(s.mux))

	s.s.Handler = s.mux
	s.s.Addr = fmt.Sprintf(":%s", HTTPPort)

	log.Fatal(s.s.ListenAndServe())
}

func (s *Server) Shutdown() {
	s.s.Shutdown(context.Background())
}

func panicHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer handlePanic(w, r)

		h.ServeHTTP(w, r)
	}
}

func handlePanic(w http.ResponseWriter, r *http.Request) {
	p := recover()
	if r == nil {
		return
	}

	var errorMessage string
	switch p.(type) {
	case error:
		errorMessage = p.(error).Error()
		break
	case string:
		errorMessage = p.(string)
		break
	default:
		errorMessage = "An internal server error occured"
		break
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(errorMessage))
}
