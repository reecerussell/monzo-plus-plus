package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/interface/http/controller"
)

// Environemt variables
var (
	HTTPPort = os.Getenv("HTTP_PORT")
)

type Server struct {
	r *mux.Router
	s *http.Server
}

func NewServer(ctn *di.Container) *Server {
	r := mux.NewRouter().StrictSlash(false)

	_ = controller.NewMonzoController(ctn, r)
	_ = controller.NewAuthController(r)
	_ = controller.NewPluginController(r)

	return &Server{
		r: r,
		s: &http.Server{},
	}
}

func (s *Server) Serve() {
	s.s.Handler = s.r
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
